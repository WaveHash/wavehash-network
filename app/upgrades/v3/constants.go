package v3

import (
	"github.com/WaveHashProtocol/wavehash/v6/app/upgrades"
	store "github.com/cosmos/cosmos-sdk/store/types"
)

const UpgradeName = "v3"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateV3UpgradeHandler,
	StoreUpgrades:        store.StoreUpgrades{},
}
