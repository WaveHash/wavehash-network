package v2

import (
	"github.com/WaveHashProtocol/wavehash/v6/app/upgrades"
	store "github.com/cosmos/cosmos-sdk/store/types"
	icacontrollertypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host/types"
)

// UpgradeName defines the on-chain upgrade name for the Wavehash v2 upgrade.
const UpgradeName = "v2" // maybe multiverse?

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV2UpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{icacontrollertypes.StoreKey, icahosttypes.StoreKey},
	},
}
