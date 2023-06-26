package v4

import (
	tokenfactorytypes "github.com/CosmWasm/token-factory/x/tokenfactory/types"
	"github.com/WaveHashProtocol/wavehash/v6/app/upgrades"
	feesharetypes "github.com/WaveHashProtocol/wavehash/v6/x/feeshare/types"
	store "github.com/cosmos/cosmos-sdk/store/types"
	icacontrollertypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/controller/types"
	ibcfeetypes "github.com/cosmos/ibc-go/v4/modules/apps/29-fee/types"
	ibchookstypes "github.com/osmosis-labs/osmosis/x/ibc-hooks/types"
	packetforwardtypes "github.com/strangelove-ventures/packet-forward-middleware/v4/router/types"
)

// UpgradeName defines the on-chain upgrade name for the upgrade.
const UpgradeName = "v4"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV4UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			tokenfactorytypes.ModuleName,
			feesharetypes.ModuleName,
			ibcfeetypes.ModuleName,
			ibchookstypes.StoreKey,
			packetforwardtypes.StoreKey,
			icacontrollertypes.StoreKey,
		},
	},
}
