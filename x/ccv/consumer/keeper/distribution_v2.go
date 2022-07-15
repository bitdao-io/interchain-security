package keeper

import (
	"bytes"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	"github.com/cosmos/interchain-security/x/ccv/utils"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"

	ibctypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	"github.com/cosmos/interchain-security/x/ccv/consumer/types"
	ccv "github.com/cosmos/interchain-security/x/ccv/types"
	gogotypes "github.com/gogo/protobuf/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const TransferTimeDelay = 1 * 7 * 24 * time.Hour // 1 weeks

// ConsumerRedistributeFrac The fraction of tokens allocated to the consumer redistribution address
// during distribution events. The fraction is a string representing a
// decimal number. For example "0.75" would represent 75%.
const ConsumerRedistributeFrac = "0.75"

func (k Keeper) Distribute(ctx sdk.Context, req abci.RequestBeginBlock) {
	// determine the total power signing the block
	var previousTotalPower, sumPreviousPrecommitPower int64
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		previousTotalPower += voteInfo.Validator.Power
		if voteInfo.SignedLastBlock {
			sumPreviousPrecommitPower += voteInfo.Validator.Power
		}
	}

	// TODO this is Tendermint-dependent
	// ref https://github.com/cosmos/cosmos-sdk/issues/3095
	if ctx.BlockHeight() > 1 {
		previousProposer := k.GetPreviousProposerConsAddr(ctx)
		if previousProposer != nil {
			k.AllocateTokens(ctx, sumPreviousPrecommitPower, previousTotalPower, previousProposer, req.LastCommitInfo.GetVotes())
		}
	}

	// record the proposer for when we payout on the next block
	consAddr := sdk.ConsAddress(req.Header.ProposerAddress)
	k.SetPreviousProposerConsAddr(ctx, consAddr)
}

func (k Keeper) AllocateTokens(
	ctx sdk.Context, sumPreviousPrecommitPower, totalPreviousPower int64,
	previousProposer sdk.ConsAddress, bondedVotes []abci.VoteInfo,
) {
	consumerFeePoolAddr := k.authKeeper.GetModuleAccount(ctx, k.feeCollectorName).GetAddress()
	fpTokens := k.bankKeeper.GetAllBalances(ctx, consumerFeePoolAddr)

	if !fpTokens.IsZero() { // don't tokens distribution if no balance in the fee collector
		if totalPreviousPower == 0 {
			if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName,
				types.ConsumerRedistributeName, fpTokens); err != nil {
				panic(err)
			}
			return
		}

		// split the fee pool, send the consumer's fraction to the consumer redistribution address
		frac, err := sdk.NewDecFromStr(ConsumerRedistributeFrac)
		if err != nil {
			panic(err)
		}
		decFPTokens := sdk.NewDecCoinsFromCoins(fpTokens...)
		consRedistrTokens, _ := decFPTokens.MulDec(frac).TruncateDecimal()
		err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName,
			types.ConsumerRedistributeName, consRedistrTokens)
		if err != nil {
			panic(err)
		}

		remainingTokens := fpTokens.Sub(consRedistrTokens)
		err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName,
			types.ProviderRewardStagingName, remainingTokens)
	}

	// calculate fraction votes
	previousFractionVotes := sdk.NewDec(sumPreviousPrecommitPower).Quo(sdk.NewDec(totalPreviousPower))

	// calculate previous proposer reward
	baseProposerReward := k.GetBaseProposerReward(ctx)
	bonusProposerReward := k.GetBonusProposerReward(ctx)
	proposerMultiplier := baseProposerReward.Add(bonusProposerReward.MulTruncate(previousFractionVotes))

	for _, vote := range bondedVotes {
		power := vote.Validator.GetPower()
		if bytes.Compare(vote.Validator.GetAddress(), previousProposer) == 0 {
			proposerPower := sdk.NewDec(totalPreviousPower).MulTruncate(proposerMultiplier)
			power += proposerPower.RoundInt64()
		}
		k.AddToProviderValidatorHoldingPool(ctx, vote.Validator.GetAddress(), power)
	}
}

// GetPreviousProposerConsAddr returns the proposer consensus address for the
// current block.
func (k Keeper) GetPreviousProposerConsAddr(ctx sdk.Context) sdk.ConsAddress {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ProposerKey)
	if bz == nil {
		return nil
	}

	addrValue := gogotypes.BytesValue{}
	k.cdc.MustUnmarshal(bz, &addrValue)
	return addrValue.GetValue()
}

// SetPreviousProposerConsAddr set the proposer public key for this block
func (k Keeper) SetPreviousProposerConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.BytesValue{Value: consAddr})
	store.Set(types.ProposerKey, bz)
}

func (k Keeper) AddToProviderValidatorHoldingPool(ctx sdk.Context, validator sdk.ConsAddress, votingPower int64) {
	holdingPool, found := k.GetProviderValidatorHoldingPool(ctx, validator)
	if !found {
		holdingPool = types.ValidatorHoldingPool{
			Weight: sdk.ZeroInt(),
		}
	}
	holdingPool.Weight = holdingPool.Weight.Add(sdk.NewInt(votingPower))
	k.SetProviderValidatorHoldingPool(ctx, validator, holdingPool)
}

// IterateValidatorHoldingPools iterates through the validator pools in the store
func (k Keeper) IterateValidatorHoldingPools(ctx sdk.Context, cb func(valAddr []byte, weight sdk.Int) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.ValidatorHoldingPoolPrefix))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		valAddr := iterator.Key()[len([]byte(types.ValidatorHoldingPoolPrefix)):]
		var validatorHoldingPool types.ValidatorHoldingPool
		k.cdc.MustUnmarshal(iterator.Value(), &validatorHoldingPool)
		if cb(valAddr, validatorHoldingPool.Weight) {
			break
		}
	}
}

func (k Keeper) SetProviderValidatorHoldingPool(ctx sdk.Context, validator sdk.ConsAddress, vhp types.ValidatorHoldingPool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&vhp)
	store.Set(types.GetValidatorHoldingPoolKey(validator), bz)
}

func (k Keeper) GetProviderValidatorHoldingPool(ctx sdk.Context, validator sdk.ConsAddress) (vhp types.ValidatorHoldingPool, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValidatorHoldingPoolKey(validator))
	if bz == nil {
		return
	}
	k.cdc.MustUnmarshal(bz, &vhp)
	return vhp, true
}

func (k Keeper) DeleteValidatorHoldingPool(ctx sdk.Context, validator sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorHoldingPoolKey(validator))
}

func (k Keeper) DistributeToProviderValidatorSetV2(ctx sdk.Context, transferTokens bool) (err error) {
	ltbh, err := k.GetLastTransmissionBlockHeight(ctx)
	if err != nil {
		return err
	}
	bpdt := k.GetBlocksPerDistributionTransmission(ctx)
	curHeight := ctx.BlockHeight()

	if (curHeight - ltbh.Height) < bpdt {
		// not enough blocks have passed for  a transmission to occur
		return nil
	}

	transmissionChannel := k.GetDistributionTransmissionChannel(ctx)
	if len(transmissionChannel) == 0 {
		return
	}

	addresses := make([][]byte, 0)
	weights := make([]sdk.Int, 0)
	totalWeights := sdk.ZeroInt()
	k.IterateValidatorHoldingPools(ctx, func(valAddr []byte, weight sdk.Int) bool {
		addresses = append(addresses, valAddr)
		weights = append(weights, weight)
		totalWeights = totalWeights.Add(weight)

		k.DeleteValidatorHoldingPool(ctx, valAddr) // clear the holding pool
		return false
	})

	providerRewardStagingAddress := k.authKeeper.GetModuleAccount(ctx, types.ProviderRewardStagingName).GetAddress()
	tokensInStagingAddress := k.bankKeeper.GetAllBalances(ctx, providerRewardStagingAddress)
	if tokensInStagingAddress.IsZero() {
		return k.resetLastTransmissionBlockHeight(ctx) // reset, return
	}

	tokens := sdk.Coins{}
	for _, token := range tokensInStagingAddress {
		if token.IsZero() {
			continue
		}
		// since SendPacket did not prefix the denomination, we must prefix denomination here
		sourcePrefix := ibctypes.GetDenomPrefix(transfertypes.PortID, transmissionChannel)
		// NOTE: sourcePrefix contains the trailing "/"
		prefixedDenom := sourcePrefix + token.GetDenom()
		// construct the denomination trace from the full raw denomination
		denomTrace := ibctypes.ParseDenomTrace(prefixedDenom)
		voucherDenom := denomTrace.IBCDenom()
		tokens = tokens.Add(sdk.NewCoin(voucherDenom, token.Amount))
	}

	providerPoolWeights := ccv.ProviderPoolWeights{
		Addresses:   addresses,
		Weights:     weights,
		TotalWeight: totalWeights,
		Tokens:      tokens,
	}

	// transfer tokens from reward staging address to toSendToProviderTokens address
	if err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ProviderRewardStagingName,
		types.ConsumerToSendToProviderName, tokensInStagingAddress); err != nil {
		return err
	}

	// transfer tokens from toSendToProviderTokens address to provider via IBC
	// empty out the toSendToProviderTokens address
	tstProviderAddr := k.authKeeper.GetModuleAccount(ctx,
		types.ConsumerToSendToProviderName).GetAddress()
	tstProviderTokens := k.bankKeeper.GetAllBalances(ctx, tstProviderAddr)
	providerAddr := k.GetProviderDistributionAddrStr(ctx)
	timeoutHeight := clienttypes.ZeroHeight()
	timeoutTimestamp := uint64(ctx.BlockTime().Add(TransferTimeDelay).UnixNano())
	if transferTokens {
		for _, token := range tstProviderTokens {
			if err = k.ibcTransferKeeper.SendTransfer(ctx,
				transfertypes.PortID,
				transmissionChannel,
				token,
				tstProviderAddr,
				providerAddr,
				timeoutHeight,
				timeoutTimestamp,
			); err != nil {
				return err
			}
		}
	}

	// if not, append provider pool weights to pending provider pool weights
	channelID, ok := k.GetProviderChannel(ctx)
	if !ok {
		k.AppendPendingProviderPoolWeights(ctx, providerPoolWeights)
		return k.resetLastTransmissionBlockHeight(ctx)
	}

	//return nil
	//send packet over IBC
	if err = utils.SendIBCPacket(
		ctx,
		k.scopedKeeper,
		k.channelKeeper,
		channelID,    // source channel id
		types.PortID, // source port id
		providerPoolWeights.GetBytes(),
	); err != nil {
		return err
	}

	return k.resetLastTransmissionBlockHeight(ctx)
}

func (k Keeper) resetLastTransmissionBlockHeight(ctx sdk.Context) error {
	newLtbh := types.LastTransmissionBlockHeight{
		Height: ctx.BlockHeight(),
	}
	return k.SetLastTransmissionBlockHeight(ctx, newLtbh)
}

func (k Keeper) GetLastTransmissionBlockHeight(ctx sdk.Context) (
	*types.LastTransmissionBlockHeight, error) {

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LastDistributionTransmissionKey)
	ltbh := &types.LastTransmissionBlockHeight{}
	if bz != nil {
		err := ltbh.Unmarshal(bz)
		if err != nil {
			return ltbh, err
		}
	}
	return ltbh, nil
}

func (k Keeper) SetLastTransmissionBlockHeight(ctx sdk.Context,
	ltbh types.LastTransmissionBlockHeight) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := ltbh.Marshal()
	if err != nil {
		return err
	}
	store.Set(types.LastDistributionTransmissionKey, bz)
	return nil
}

func (k Keeper) ChannelOpenInit(ctx sdk.Context, msg *channeltypes.MsgChannelOpenInit) (
	*channeltypes.MsgChannelOpenInitResponse, error) {
	return k.ibcCoreKeeper.ChannelOpenInit(sdk.WrapSDKContext(ctx), msg)
}

func (k Keeper) GetConnectionHops(ctx sdk.Context, srcPort, srcChan string) ([]string, error) {
	ch, found := k.channelKeeper.GetChannel(ctx, srcPort, srcChan)
	if !found {
		return []string{}, sdkerrors.Wrapf(ccv.ErrChannelNotFound,
			"cannot get connection hops from non-existent channel")
	}
	return ch.ConnectionHops, nil
}
