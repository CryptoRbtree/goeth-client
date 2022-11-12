package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func block() {
	var (
		ctx         = context.Background()
		url         = "wss://eth-mainnet.g.alchemy.com/v2/" + os.Getenv("ALCHEMY_ID")
		client, err = ethclient.DialContext(ctx, url)
	)
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("latest block number =", header.Number.String())

	blockNumber := big.NewInt(15500000)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("block number =", block.Number().Uint64())
	fmt.Println("block timestamp =", block.Time())
	fmt.Println("block difficulty =", block.Difficulty().Uint64())
	fmt.Println("block hash =", block.Hash().Hex())
	fmt.Println("block transaction count =", len(block.Transactions()))

	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("block transaction count =", count)

	blockSubscribe(client)
}

func blockSubscribe(client *ethclient.Client) {
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	prevBaseFee := big.NewFloat(0)

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			fmt.Println("----new block mined----")

			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("block number:", block.Number().Uint64())
			fmt.Println("block hash:", block.Hash().Hex())
			fmt.Println("block timestamp:", block.Time())
			currBaseFee := big.NewFloat(0)
			currBaseFee.SetString(block.BaseFee().String())
			currBaseFee.Quo(currBaseFee, big.NewFloat(1e9))
			increaseRate := big.NewFloat(0)
			increaseRate.Quo(currBaseFee, prevBaseFee)
			increaseRate.Add(increaseRate, big.NewFloat(-1))
			increaseRate.Mul(increaseRate, big.NewFloat(100))
			fmt.Println("block basefee(Gwei)/increase_rate(%):", currBaseFee, increaseRate)
			prevBaseFee = currBaseFee
			fmt.Println("block gas used:", block.GasUsed())
			fmt.Println("block transaction count:", len(block.Transactions()))
		}
	}
}
