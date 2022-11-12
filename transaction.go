package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func transaction() {
	var (
		ctx         = context.Background()
		url         = "wss://eth-mainnet.g.alchemy.com/v2/" + os.Getenv("ALCHEMY_ID")
		client, err = ethclient.DialContext(ctx, url)
	)
	if err != nil {
		log.Fatal(err)
	}

	// 1
	blockNumber := big.NewInt(15000023)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	transactions := block.Transactions()
	fmt.Println("total transactions count:", len(transactions))
	for ind, tx := range transactions {
		fmt.Println("----", ind, "----")
		fmt.Println("transaction hash:", tx.Hash().Hex())
		fmt.Println("transaction value:", tx.Value().String())
		fmt.Println("transaction gas limit:", tx.Gas())
		fmt.Println("transaction fee cap per gas:", tx.GasFeeCap())
		fmt.Println("transaction tip cap per gas:", tx.GasTipCap())
		fmt.Println("transaction gas price:", tx.GasPrice())
		fmt.Println("transaction nonce:", tx.Nonce())                   // 110644
		fmt.Println("transaction data:", hex.EncodeToString(tx.Data())) // []
		fmt.Println("transaction to:", tx.To().Hex())                   // 0x55fE59D8Ad77035154dDd0AD0388D09Dd4047A8e

		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		baseFee := big.NewInt(1000000000)
		if msg, err := tx.AsMessage(types.NewEIP155Signer(chainID), baseFee); err == nil {
			fmt.Println("transaction from:", msg.From().Hex())
		}

		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("transaction gas used:", receipt.GasUsed)
		fmt.Println("transaction status:", receipt.Status)
	}

	// 2
	blockHash := common.HexToHash("0x92a025d0798bc1bf1b284fa1440007a2c0991bb65e0bb6c72bf9d6a4c387b195") // 15945766
	count, err := client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("----block 15945766----")
	for idx := uint(0); idx < count; idx++ {
		tx, err := client.TransactionInBlock(context.Background(), blockHash, idx)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("transaction hash:", tx.Hash().Hex())
	}

	// 3
	fmt.Println("----last ethw transaction----")
	txHash := common.HexToHash("0xec9db5bfbcd30ad2e3070b626ed4f78abce88687c5d1eb23464242be5edcb537")
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("transaction hash:", tx.Hash().Hex())
	fmt.Println("transaction gas limit:", tx.Gas())
	fmt.Println("isPending:", isPending) // false
}
