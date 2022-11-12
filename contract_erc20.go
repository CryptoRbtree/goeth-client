package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	token "github.com/CryptoRbtree/goeth-client/contract-erc20"
)

func contractErc20() {
	var (
		ctx         = context.Background()
		url         = "https://eth-mainnet.g.alchemy.com/v2/" + os.Getenv("ALCHEMY_ID")
		client, err = ethclient.DialContext(ctx, url)
	)
	if err != nil {
		log.Fatal(err)
	}

	tokenAddress := common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F") // DAI Address
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress("0x245cc372c84b3645bf0ffe6538620b04a217988b") // Olympus funds
	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err)
	}

	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("name: %s\n", name)
	fmt.Printf("symbol: %s\n", symbol)
	fmt.Printf("decimals: %v\n", decimals)

	fbal := new(big.Float)
	fbal.SetString(bal.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))
	fmt.Printf("balance: %f\n", value)
}
