package main

import (
	"os"

	"github.com/WaveHashProtocol/wavehash/v6/app"
	"github.com/WaveHashProtocol/wavehash/v6/cmd/wavehashd/cmd"
	"github.com/cosmos/cosmos-sdk/server"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	app.SetAddressPrefixes()
	rootCmd, _ := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
