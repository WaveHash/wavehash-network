package app

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
)

const (
	// DefaultWavehashInstanceCost is initially set the same as in wasmd
	DefaultWavehashInstanceCost uint64 = 60_000
	// DefaultWavehashCompileCost set to a large number for testing
	DefaultWavehashCompileCost uint64 = 3
)

// WavehashGasRegisterConfig is defaults plus a custom compile amount
func WavehashGasRegisterConfig() wasmkeeper.WasmGasRegisterConfig {
	gasConfig := wasmkeeper.DefaultGasRegisterConfig()
	gasConfig.InstanceCost = DefaultWavehashInstanceCost
	gasConfig.CompileCost = DefaultWavehashCompileCost

	return gasConfig
}

func NewWavehashWasmGasRegister() wasmkeeper.WasmGasRegister {
	return wasmkeeper.NewWasmGasRegister(WavehashGasRegisterConfig())
}
