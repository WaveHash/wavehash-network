package keeper_test

import (
	"encoding/json"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	wavehashapp "github.com/WaveHashProtocol/wavehash/v6/app"
	"github.com/cosmos/cosmos-sdk/simapp"

	"github.com/WaveHashProtocol/wavehash/v6/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// returns context and an app with updated mint keeper
func createTestApp(isCheckTx bool) (*wavehashapp.App, sdk.Context) { //nolint:unparam
	app := setup(isCheckTx)

	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.MintKeeper.SetParams(ctx, types.DefaultParams())
	app.MintKeeper.SetMinter(ctx, types.DefaultInitialMinter())

	return app, ctx
}

func setup(isCheckTx bool) *wavehashapp.App {
	app, genesisState := genApp(!isCheckTx, 5)
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

func genApp(withGenesis bool, invCheckPeriod uint) (*wavehashapp.App, wavehashapp.GenesisState) {
	db := dbm.NewMemDB()
	encCdc := wavehashapp.MakeEncodingConfig()
	app := wavehashapp.New(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		simapp.DefaultNodeHome,
		invCheckPeriod,
		encCdc,
		wavehashapp.GetEnabledProposals(),
		simapp.EmptyAppOptions{},
		wavehashapp.GetWasmOpts(simapp.EmptyAppOptions{}),
	)

	if withGenesis {
		return app, wavehashapp.NewDefaultGenesisState(encCdc.Marshaler)
	}

	return app, wavehashapp.GenesisState{}
}
