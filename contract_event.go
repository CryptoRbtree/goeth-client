package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	token "github.com/CryptoRbtree/goeth-client/contract-erc20"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// LogTransfer ..
type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}

func contractEvent() {
	var (
		ctx         = context.Background()
		url         = "wss://eth-mainnet.g.alchemy.com/v2/" + os.Getenv("ALCHEMY_ID")
		client, err = ethclient.DialContext(ctx, url)
	)
	if err != nil {
		log.Fatal(err)
	}

	tokenAddress := common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F") // DAI Address
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(15500000),
		ToBlock:   big.NewInt(15500005),
		Addresses: []common.Address{
			tokenAddress,
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(token.TokenABI)))
	if err != nil {
		log.Fatal(err)
	}

	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

	for _, vLog := range logs {
		fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
		fmt.Printf("Log Index: %d\n", vLog.Index)

		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			fmt.Printf("Log Name: Transfer\n")

			var transferEvent LogTransfer
			err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
			fmt.Printf("From: %s\n", transferEvent.From.Hex())
			fmt.Printf("To: %s\n", transferEvent.To.Hex())
			fmt.Printf("Value: %s\n", transferEvent.Tokens.String())
		}

		fmt.Printf("\n")
	}

	eventSubscribe(client)
}

func eventSubscribe(client *ethclient.Client) {
	contractAddress := common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("----subscribe DAI events----\n\n")
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Printf("%+v\n\n", vLog) // pointer to log event
		}
	}
}
