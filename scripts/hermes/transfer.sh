#!/bin/sh

CHAIN_A_ARGS="--from wavehash1 --keyring-backend test --chain-id local-1 --home $HOME/.wavehash1/ --node http://localhost:26657 --yes"

# Send from local-1 to local-2 via the relayer
wavehashd tx ibc-transfer transfer transfer channel-0 wavehash1hj5fveer5cjtn4wd6wstzugjfdxzl0xps73ftl 9uwaha $CHAIN_A_ARGS --packet-timeout-height 0-0

sleep 6

# check the query on the other chain to ensure it went through
wavehashd q bank balances wavehash1hj5fveer5cjtn4wd6wstzugjfdxzl0xps73ftl --chain-id local-2 --node http://localhost:36657