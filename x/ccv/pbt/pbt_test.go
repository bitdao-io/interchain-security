package pbt_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/v3/modules/core/23-commitment/types"
	ibctmtypes "github.com/cosmos/ibc-go/v3/modules/light-clients/07-tendermint/types"
	ibctesting "github.com/cosmos/ibc-go/v3/testing"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	appConsumer "github.com/cosmos/interchain-security/app/consumer"
	appProvider "github.com/cosmos/interchain-security/app/provider"
	"github.com/cosmos/interchain-security/testutil/simapp"
	consumertypes "github.com/cosmos/interchain-security/x/ccv/consumer/types"
	providertypes "github.com/cosmos/interchain-security/x/ccv/provider/types"
	"github.com/cosmos/interchain-security/x/ccv/types"

	tmtypes "github.com/tendermint/tendermint/types"
)

type PBTTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains
	providerChain *ibctesting.TestChain
	consumerChain *ibctesting.TestChain

	path *ibctesting.Path
}

const P = "provider"
const C = "consumer"

// TODO: do I need different denoms for each chain?
const denom = sdk.DefaultBondDenom
const maxValidators = 2

func init() {
	// Tokens = Power
	sdk.DefaultPowerReduction = sdk.NewInt(1)
}

func TestPBTTestSuite(t *testing.T) {
	suite.Run(t, new(PBTTestSuite))
}

func (s *PBTTestSuite) SetupTest() {

	s.coordinator, s.providerChain, s.consumerChain = simapp.NewProviderConsumerCoordinator(s.T())

	// s.DisableConsumerDistribution() TODO: needed??

	tmConfig := ibctesting.NewTendermintConfig()

	// commit a block on provider chain before creating client
	s.coordinator.CommitBlock(s.providerChain)

	// create client and consensus state of provider chain to initialize consumer chain genesis.
	height := s.providerChain.LastHeader.GetHeight().(clienttypes.Height)
	UpgradePath := []string{"upgrade", "upgradedIBCState"}

	providerClient := ibctmtypes.NewClientState(
		s.providerChain.ChainID, tmConfig.TrustLevel, tmConfig.TrustingPeriod, tmConfig.UnbondingPeriod, tmConfig.MaxClockDrift,
		height, commitmenttypes.GetSDKSpecs(), UpgradePath, tmConfig.AllowUpdateAfterExpiry, tmConfig.AllowUpdateAfterMisbehaviour,
	)
	providerConsState := s.providerChain.LastHeader.ConsensusState()

	valUpdates := tmtypes.TM2PB.ValidatorUpdates(s.providerChain.Vals)

	params := consumertypes.NewParams(
		true,
		1000, // about 2 hr at 7.6 seconds per blocks
		"",
		"",
		"0.5", // 50%
	)
	consumerGenesis := consumertypes.NewInitialGenesisState(providerClient, providerConsState, valUpdates, params)
	s.consumerChain.App.(*appConsumer.App).ConsumerKeeper.InitGenesis(s.ctx(C), consumerGenesis)

	s.path = ibctesting.NewPath(s.consumerChain, s.providerChain)
	s.path.EndpointA.ChannelConfig.PortID = consumertypes.PortID
	s.path.EndpointB.ChannelConfig.PortID = providertypes.PortID
	s.path.EndpointA.ChannelConfig.Version = types.Version
	s.path.EndpointB.ChannelConfig.Version = types.Version
	s.path.EndpointA.ChannelConfig.Order = channeltypes.ORDERED
	s.path.EndpointB.ChannelConfig.Order = channeltypes.ORDERED

	providerClientId, ok := s.consumerChain.App.(*appConsumer.App).ConsumerKeeper.GetProviderClient(s.ctx(C))
	if !ok {
		panic("must already have provider client on consumer chain")
	}

	// set consumer endpoint's clientID
	s.path.EndpointA.ClientID = providerClientId

	// TODO: No idea why or how this works, but it seems that it needs to be done.
	s.path.EndpointB.Chain.SenderAccount.SetAccountNumber(6)
	s.path.EndpointA.Chain.SenderAccount.SetAccountNumber(6)

	// create consumer client on provider chain and set as consumer client for consumer chainID in provider keeper.
	s.path.EndpointB.CreateClient()
	s.providerChain.App.(*appProvider.App).ProviderKeeper.SetConsumerClient(s.ctx(P), s.consumerChain.ChainID, s.path.EndpointB.ClientID)

	// TODO: I added this section, should I remove it or move it?
	//~~~~~~~~~~
	s.coordinator.CreateConnections(s.path)

	// CCV channel handshake will automatically initiate transfer channel handshake on ACK
	// so transfer channel will be on stage INIT when CreateChannels for ccv path returns.
	s.coordinator.CreateChannels(s.path)
	//~~~~~~~~~~
}

/*
~~~~~~~~~~~~
QUERIES
~~~~~~~~~~~~
*/

func (s *PBTTestSuite) ctx(chain string) sdk.Context {
	return s.chain(chain).GetContext()
}

func (s *PBTTestSuite) chain(chain string) *ibctesting.TestChain {
	chains := make(map[string]*ibctesting.TestChain)
	chains["provider"] = s.providerChain
	chains["consumer"] = s.consumerChain
	return chains[chain]
}

func (s *PBTTestSuite) height(chain string) int64 {
	return s.chain(chain).CurrentHeader.GetHeight()
}

func (s *PBTTestSuite) endpoint(chain string) *ibctesting.Endpoint {
	endpoints := make(map[string]*ibctesting.Endpoint)
	endpoints["provider"] = s.path.EndpointB
	endpoints["consumer"] = s.path.EndpointA
	return endpoints[chain]
}

func (s *PBTTestSuite) delegator() sdk.AccAddress {
	delAddr := s.providerChain.SenderAccount.GetAddress()
	return delAddr
}

func (s *PBTTestSuite) validator(i int64) sdk.ValAddress {
	tmValidator := s.providerChain.Vals.Validators[i]
	valAddr, err := sdk.ValAddressFromHex(tmValidator.Address.String())
	s.Require().NoError(err)
	return valAddr
}

func (s *PBTTestSuite) consAddr(i int64) sdk.ConsAddress {
	val := s.providerChain.Vals.Validators[i]
	consAddr := sdk.ConsAddress(val.Address)
	return consAddr
}

func (s *PBTTestSuite) validatorStatus(chain string, i int64) stakingtypes.BondStatus {
	addr := s.validator(i)
	val, found := s.chain(chain).App.GetStakingKeeper().GetValidator(s.ctx(chain), addr)
	if !found {
		s.T().Fatal("Couldn't GetValidator")
	}
	return val.GetStatus()
}

func (s *PBTTestSuite) delegatorBalance() int64 {
	del := s.delegator()
	app := s.providerChain.App.(*appProvider.App)
	bal := app.BankKeeper.GetBalance(s.ctx(P), del, denom)
	return bal.Amount.Int64()
}

func (s *PBTTestSuite) validatorTokens(chain string, i int64) int64 {
	addr := s.validator(i)
	val, found := s.chain(chain).App.GetStakingKeeper().GetValidator(s.ctx(chain), addr)
	if !found {
		s.T().Fatal("Couldn't GetValidator")
	}
	return val.Tokens.Int64()
}

func (s *PBTTestSuite) delegation(i int64) int64 {
	addr := s.delegator()
	del, found := s.providerChain.App.GetStakingKeeper().GetDelegation(s.ctx(P), addr, s.validator(i))
	if !found {
		s.T().Fatal("Couldn't GetDelegation")
	}
	return del.Shares.TruncateInt64()
}

func equalHeights(s *PBTTestSuite) {
	ph := s.height(P)
	ch := s.height(C)
	if ph != ch {
		s.T().Fatal("Bad test")
	}
}

/*
~~~~~~~~~~~~
MODEL
~~~~~~~~~~~~
*/

type Action struct {
	Amt             int      `json:"amt,omitempty"`
	Chain           string   `json:"chain,omitempty"`
	Chains          []string `json:"chains,omitempty"`
	Factor          float64  `json:"factor,omitempty"`
	Height          int      `json:"height,omitempty"`
	IsDowntime      bool     `json:"is_downtime,omitempty"`
	Kind            string   `json:"kind"`
	N               int      `json:"n,omitempty"`
	Power           int      `json:"power,omitempty"`
	SecondsPerBlock int      `json:"seconds_per_block,omitempty"`
	Val             int      `json:"val,omitempty"`
}

type Trace struct {
	Actions []Action `json:"actions"`
}

type Delegate struct {
	val int64
	amt int64
}
type Undelegate struct {
	val int64
	amt int64
}
type JumpNBlocks struct {
	chains          []string
	n               int64
	secondsPerBlock int64
}
type Deliver struct {
	chain string
}
type ProviderSlash struct {
	val    int64
	power  int64
	height int64
	factor int64
}
type ConsumerSlash struct {
	val        int64
	height     int64
	power      int64
	isDowntime bool
}

func (s *PBTTestSuite) delegate(a Delegate) {
	psk := s.providerChain.App.GetStakingKeeper()
	pskServer := stakingkeeper.NewMsgServerImpl(psk)
	amt := sdk.NewCoin(denom, sdk.NewInt(a.amt))
	del := s.delegator()
	val := s.validator(a.val)
	msg := stakingtypes.NewMsgDelegate(del, val, amt)
	pskServer.Delegate(sdk.WrapSDKContext(s.ctx(P)), msg)
}

func (s *PBTTestSuite) undelegate(a Undelegate) {

	psk := s.providerChain.App.GetStakingKeeper()
	pskServer := stakingkeeper.NewMsgServerImpl(psk)
	amt := sdk.NewCoin(denom, sdk.NewInt(a.amt))
	del := s.delegator()
	val := s.validator(a.val)
	msg := stakingtypes.NewMsgUndelegate(del, val, amt)
	pskServer.Undelegate(sdk.WrapSDKContext(s.ctx(P)), msg)
}

func (s *PBTTestSuite) endBlock(chain string) {

}
func (s *PBTTestSuite) increaseSeconds(seconds int64) {
	s.coordinator.IncrementTimeBy(time.Second * time.Duration(seconds))
}
func (s *PBTTestSuite) deliver(a Deliver) {
	// TODO:!
}

func (s *PBTTestSuite) jumpNBlocks(a JumpNBlocks) {
	for i := int64(0); i < a.n; i++ {
		for _, c := range a.chains {
			s.endBlock(c)
		}
		s.increaseSeconds(a.secondsPerBlock)
	}
}

func (s *PBTTestSuite) providerSlash(a ProviderSlash) {
	psk := s.providerChain.App.GetStakingKeeper()
	val := s.consAddr(a.val)
	h := int64(a.height)
	power := int64(a.power)
	factor := sdk.NewDec(int64(a.factor))
	psk.Slash(s.ctx(P), val, h, power, factor, -1) // TODO: check params here!
}

func (s *PBTTestSuite) consumerSlash(a ConsumerSlash) {
	cccvk := s.consumerChain.App.(*appConsumer.App).ConsumerKeeper
	val := s.consAddr(a.val)
	h := int64(a.height)
	power := int64(a.power)
	kind := stakingtypes.Downtime
	if !a.isDowntime {
		kind = stakingtypes.DoubleSign
	}
	cccvk.Slash(s.ctx(C), val, h, power, sdk.Dec{}, kind) // TODO: check params here!
}

/*
~~~~~~~~~~~~
ZERO TEST
~~~~~~~~~~~~
*/

// TODO: clear up these hacks after stripping provider/consumer
func (s *PBTTestSuite) DisableConsumerDistribution() {
	cChain := s.consumerChain
	cApp := cChain.App.(*appConsumer.App)
	for i, moduleName := range cApp.MM.OrderBeginBlockers {
		if moduleName == distrtypes.ModuleName {
			cApp.MM.OrderBeginBlockers = append(cApp.MM.OrderBeginBlockers[:i], cApp.MM.OrderBeginBlockers[i+1:]...)
			return
		}
	}
}

func adjustParams(s *PBTTestSuite) {
	params := s.providerChain.App.GetStakingKeeper().GetParams(s.ctx(P))
	params.MaxValidators = maxValidators
	s.providerChain.App.GetStakingKeeper().SetParams(s.ctx(P), params)
}

func (s *PBTTestSuite) TestAssumptions() {

	adjustParams(s)

	s.jumpNBlocks(JumpNBlocks{[]string{P}, 1, 5})
	// TODO: Is it correct to catch the consumer up with the provider here?
	s.jumpNBlocks(JumpNBlocks{[]string{C}, 2, 5})

	equalHeights(s)

	/*
		delegatorBalance() overflows int64 because it is set to a number greater than 2^63 in genesis.
		It's easiest to assume we have enough funds.
	*/

	maxValsE := uint32(3)
	maxVals := s.providerChain.App.GetStakingKeeper().GetParams(s.ctx(P)).MaxValidators

	if maxValsE != maxVals {
		s.T().Fatal("Bad test")
	}

	for i := 0; i < 4; i++ {
		// This is the genesis delegation
		delE := int64(1)
		del := s.delegation(int64(i))
		if delE != del {
			s.T().Fatal("Bad test")
		}
	}

	step := int64(1)

	for i := 0; i < 3; i++ {
		s.delegate(Delegate{int64(i), (3 - int64(i)) * step})
	}

	for i := 0; i < 4; i++ {
		delE := (3-int64(i))*step + int64(1)
		del := s.delegation(int64(i))
		if delE != del {
			s.T().Fatal("Bad test")
		}
	}

	for i := 0; i < maxValidators; i++ {
		// First ones should be bonded
		A := s.validatorStatus(P, int64(i))
		E := stakingtypes.Bonded
		if E != A {
			s.T().Fatal("Bad test")
		}
	}

	for i := maxValidators; i < 4; i++ {
		A := s.validatorStatus(P, int64(i))
		// Last one is unbonding
		E := stakingtypes.Unbonding
		if E != A {
			s.T().Fatal("Bad test")
		}
	}

	equalHeights(s)

	for i := 0; i < maxValidators; i++ {
		A := s.validatorTokens(P, int64(i))
		E := (3-int64(i))*step + int64(1)
		if E != A {
			s.T().Fatal("Bad test")
		}
	}

	// validator tokens are 4,3,2,1

}

/*
~~~~~~~~~~~~
TRACE TEST
~~~~~~~~~~~~
*/

func executeTrace(s *PBTTestSuite, trace Trace) {

	for _, a := range trace.Actions {
		switch a.Kind {
		case "delegate":
			s.delegate(Delegate{
				int64(a.Val),
				int64(a.Amt),
			})
		case "undelegate":
			s.undelegate(Undelegate{
				int64(a.Val),
				int64(a.Amt),
			})
		case "jumpNBlocks":
			s.jumpNBlocks(JumpNBlocks{
				a.Chains,
				int64(a.N),
				int64(a.SecondsPerBlock),
			})
		case "deliver":
			s.deliver(Deliver{a.Chain})
		case "providerSlash":
			s.providerSlash(ProviderSlash{
				int64(a.Val),
				int64(a.Power),
				int64(a.Height),
				1, // TODO: unhard code, use factor!
			})
		case "consumerSlash":
			s.consumerSlash(ConsumerSlash{
				int64(a.Val),
				int64(a.Height),
				int64(a.Power),
				a.IsDowntime,
			})
		}
	}
}

func loadTraces(fn string) []Trace {

	fd, err := os.Open(fn)

	if err != nil {
		panic(err)
	}

	defer fd.Close()

	byteValue, _ := ioutil.ReadAll(fd)

	var ret []Trace

	err = json.Unmarshal([]byte(byteValue), &ret)

	if err != nil {
		panic(err)
	}

	return ret
}

func executeTraces(s *PBTTestSuite, traces []Trace) {
	for _, trace := range traces {
		executeTrace(s, trace)
	}
}

func (s *PBTTestSuite) TestTracesHC() {
	executeTraces(s, loadTraces("trace/hc.json"))
}
