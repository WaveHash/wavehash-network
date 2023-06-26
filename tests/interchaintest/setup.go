package interchaintest

import (
	feesharetypes "github.com/WaveHashProtocol/wavehash/v6/x/feeshare/types"

	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/strangelove-ventures/interchaintest/v4/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v4/ibc"
)

var (
	WavehashE2ERepo  = "ghcr.io/wavehashprotocol/wavehash-e2e"
	WavehashMainRepo = "ghcr.io/wavehashprotocol/wavehash"

	wavehashRepo, wavehashVersion = GetDockerImageInfo()

	WavehashImage = ibc.DockerImage{
		Repository: wavehashRepo,
		Version:    wavehashVersion,
		UidGid:     "1025:1025",
	}

	wavehashConfig = ibc.ChainConfig{
		Type:                "cosmos",
		Name:                "wavehash",
		ChainID:             "wavehash-2",
		Images:              []ibc.DockerImage{WavehashImage},
		Bin:                 "wavehashd",
		Bech32Prefix:        "wavehash",
		Denom:               "uwaha",
		CoinType:            "118",
		GasPrices:           "0.0uwaha",
		GasAdjustment:       1.1,
		TrustingPeriod:      "112h",
		NoHostMount:         false,
		SkipGenTx:           false,
		PreGenesis:          nil,
		ModifyGenesis:       nil,
		ConfigFileOverrides: nil,
		EncodingConfig:      wavehashEncoding(),
	}

	pathWavehashGaia    = "wavehash-gaia"
	genesisWalletAmount = int64(10_000_000)
)

// wavehashEncoding registers the Wavehash specific module codecs so that the associated types and msgs
// will be supported when writing to the blocksdb sqlite database.
func wavehashEncoding() *simappparams.EncodingConfig {
	cfg := cosmos.DefaultEncoding()

	// register custom types
	feesharetypes.RegisterInterfaces(cfg.InterfaceRegistry)

	return &cfg
}
