package interchaintest

import (
	"context"
	"fmt"
	"testing"

	transfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	"github.com/strangelove-ventures/interchaintest/v4"
	"github.com/strangelove-ventures/interchaintest/v4/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v4/ibc"
	"github.com/strangelove-ventures/interchaintest/v4/testreporter"
	"github.com/strangelove-ventures/interchaintest/v4/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// TestWavehashGaiaIBCTransfer spins up a Wavehash and Gaia network, initializes an IBC connection between them,
// and sends an ICS20 token transfer from Wavehash->Gaia and then back from Gaia->Wavehash.
func TestWavehashGaiaIBCTransfer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	// Create chain factory with Wavehash and Gaia
	numVals := 1
	numFullNodes := 1

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:          "wavehash",
			ChainConfig:   wavehashConfig,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
		{
			Name:          "gaia",
			Version:       "v9.0.0",
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
	})

	// Get chains from the chain factory
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	wavehash, gaia := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

	// Create relayer factory to utilize the go-relayer
	client, network := interchaintest.DockerSetup(t)

	r := interchaintest.NewBuiltinRelayerFactory(ibc.CosmosRly, zaptest.NewLogger(t)).Build(t, client, network)

	// Create a new Interchain object which describes the chains, relayers, and IBC connections we want to use
	ic := interchaintest.NewInterchain().
		AddChain(wavehash).
		AddChain(gaia).
		AddRelayer(r, "rly").
		AddLink(interchaintest.InterchainLink{
			Chain1:  wavehash,
			Chain2:  gaia,
			Relayer: r,
			Path:    pathWavehashGaia,
		})

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	ctx := context.Background()

	err = ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:         t.Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: false,

		// This can be used to write to the block database which will index all block data e.g. txs, msgs, events, etc.
		// BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = ic.Close()
	})

	// Start the relayer
	require.NoError(t, r.StartRelayer(ctx, eRep, pathWavehashGaia))
	t.Cleanup(
		func() {
			err := r.StopRelayer(ctx, eRep)
			if err != nil {
				panic(fmt.Errorf("an error occurred while stopping the relayer: %s", err))
			}
		},
	)

	// Create some user accounts on both chains
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), genesisWalletAmount, wavehash, gaia)

	// Wait a few blocks for relayer to start and for user accounts to be created
	err = testutil.WaitForBlocks(ctx, 5, wavehash, gaia)
	require.NoError(t, err)

	// Get our Bech32 encoded user addresses
	wavehashUser, gaiaUser := users[0], users[1]

	wavehashUserAddr := wavehashUser.Bech32Address(wavehash.Config().Bech32Prefix)
	gaiaUserAddr := gaiaUser.Bech32Address(gaia.Config().Bech32Prefix)

	// Get original account balances
	wavehashOrigBal, err := wavehash.GetBalance(ctx, wavehashUserAddr, wavehash.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, genesisWalletAmount, wavehashOrigBal)

	gaiaOrigBal, err := gaia.GetBalance(ctx, gaiaUserAddr, gaia.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, genesisWalletAmount, gaiaOrigBal)

	// Compose an IBC transfer and send from Wavehash -> Gaia
	const transferAmount = int64(1_000)
	transfer := ibc.WalletAmount{
		Address: gaiaUserAddr,
		Denom:   wavehash.Config().Denom,
		Amount:  transferAmount,
	}

	channel, err := ibc.GetTransferChannel(ctx, r, eRep, wavehash.Config().ChainID, gaia.Config().ChainID)
	require.NoError(t, err)

	transferTx, err := wavehash.SendIBCTransfer(ctx, channel.ChannelID, wavehashUserAddr, transfer, ibc.TransferOptions{})
	require.NoError(t, err)

	wavehashHeight, err := wavehash.Height(ctx)
	require.NoError(t, err)

	// Poll for the ack to know the transfer was successful
	_, err = testutil.PollForAck(ctx, wavehash, wavehashHeight, wavehashHeight+10, transferTx.Packet)
	require.NoError(t, err)

	// Get the IBC denom for uwaha on Gaia
	wavehashTokenDenom := transfertypes.GetPrefixedDenom(channel.Counterparty.PortID, channel.Counterparty.ChannelID, wavehash.Config().Denom)
	wavehashIBCDenom := transfertypes.ParseDenomTrace(wavehashTokenDenom).IBCDenom()

	// Assert that the funds are no longer present in user acc on Wavehash and are in the user acc on Gaia
	wavehashUpdateBal, err := wavehash.GetBalance(ctx, wavehashUserAddr, wavehash.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, wavehashOrigBal-transferAmount, wavehashUpdateBal)

	gaiaUpdateBal, err := gaia.GetBalance(ctx, gaiaUserAddr, wavehashIBCDenom)
	require.NoError(t, err)
	require.Equal(t, transferAmount, gaiaUpdateBal)

	// Compose an IBC transfer and send from Gaia -> Wavehash
	transfer = ibc.WalletAmount{
		Address: wavehashUserAddr,
		Denom:   wavehashIBCDenom,
		Amount:  transferAmount,
	}

	transferTx, err = gaia.SendIBCTransfer(ctx, channel.Counterparty.ChannelID, gaiaUserAddr, transfer, ibc.TransferOptions{})
	require.NoError(t, err)

	gaiaHeight, err := gaia.Height(ctx)
	require.NoError(t, err)

	// Poll for the ack to know the transfer was successful
	_, err = testutil.PollForAck(ctx, gaia, gaiaHeight, gaiaHeight+10, transferTx.Packet)
	require.NoError(t, err)

	// Assert that the funds are now back on Wavehash and not on Gaia
	wavehashUpdateBal, err = wavehash.GetBalance(ctx, wavehashUserAddr, wavehash.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, wavehashOrigBal, wavehashUpdateBal)

	gaiaUpdateBal, err = gaia.GetBalance(ctx, gaiaUserAddr, wavehashIBCDenom)
	require.NoError(t, err)
	require.Equal(t, int64(0), gaiaUpdateBal)
}
