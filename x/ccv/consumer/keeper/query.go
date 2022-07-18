package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/interchain-security/x/ccv/consumer/types"
)

func (k Keeper) QueryCurrentProviderHoldingPool(c context.Context, _ *types.QueryCurrentProviderHoldingPoolRequest) (*types.CurrentProviderHoldingPool, error) {
	ctx := sdk.UnwrapSDKContext(c)
	validatorWeights := make([]*types.ValidatorWeight, 0)
	var err error
	k.IterateValidatorHoldingPools(ctx, func(valAddr []byte, weight sdk.Int) bool {
		valAddrStr, err := sdk.Bech32ifyAddressBytes(sdk.GetConfig().GetBech32ConsensusAddrPrefix(), valAddr)
		if err != nil {
			return true
		}
		vw := &types.ValidatorWeight{
			Address: valAddrStr,
			Weight:  weight,
		}
		validatorWeights = append(validatorWeights, vw)
		return false
	})

	if err != nil {
		return nil, err
	}

	ltbh, err := k.GetLastTransmissionBlockHeight(ctx)
	if err != nil {
		return nil, err
	}

	providerRewardStagingAddress := k.authKeeper.GetModuleAccount(ctx, types.ProviderRewardStagingName).GetAddress()
	tokens := k.bankKeeper.GetAllBalances(ctx, providerRewardStagingAddress)

	return &types.CurrentProviderHoldingPool{
		ValidatorWeights: validatorWeights,
		StartHeight:      ltbh.GetHeight() + 1,
		CurrentHeight:    ctx.BlockHeight(),
		Tokens:           tokens,
	}, nil
}
