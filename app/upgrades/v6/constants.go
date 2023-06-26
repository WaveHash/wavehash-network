package v6

import (
	"github.com/WaveHashProtocol/wavehash/v6/app/upgrades"
	store "github.com/cosmos/cosmos-sdk/store/types"
)

// UpgradeName defines the on-chain upgrade name for the upgrade.
const UpgradeName = "v6"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV6PatchUpgradeHandler,
	StoreUpgrades:        store.StoreUpgrades{},
}
