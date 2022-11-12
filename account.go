package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func account() {
	var (
		ctx         = context.Background()
		url         = "https://eth-mainnet.g.alchemy.com/v2/" + os.Getenv("ALCHEMY_ID")
		client, err = ethclient.DialContext(ctx, url)
	)

	if err != nil {
		log.Fatal(err)
	}

	// EOA
	account := common.HexToAddress("0xab5801a7d398351b8be11c439e05c5b3259aec9b") // vitalik
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("account current eth balance =", balance, "wei")

	blockNumber := big.NewInt(15500000)
	balanceAt, err := client.BalanceAt(context.Background(), account, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	balanceAtFloat := big.NewFloat(0)
	balanceAtFloat.SetString(balanceAt.String())
	balanceAtFloat.Quo(balanceAtFloat, big.NewFloat(1e18))
	fmt.Println("account blocknum 15500000 eth balance =", balanceAtFloat, "eth")

	// CA
	contractAccount := common.HexToAddress("0x220866B1A2219f40e72f5c628B65D54268cA3A9D") // vitalik multisig
	contractBalance, err := client.BalanceAt(context.Background(), contractAccount, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("smart contract account eth balance =", contractBalance, "wei")
	code, err := client.CodeAt(context.Background(), contractAccount, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("smart contract deployed code =", hex.EncodeToString(code))
}
