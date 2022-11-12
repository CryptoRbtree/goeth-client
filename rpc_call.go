package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
)

// https://stackoverflow.com/questions/53237759/how-to-correctly-send-rpc-call-using-golang-to-get-smart-contract-owner/53260846#53260846

func rpcCall() {
	var (
		url         = "https://eth-mainnet.g.alchemy.com/v2/" + os.Getenv("ALCHEMY_ID")
		client, err = rpc.DialHTTP(url)
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	type request struct {
		To   string `json:"to"`
		Data string `json:"data"`
	}

	usdt := "0xdAC17F958D2ee523a2206206994597C13D831ec7"
	functionHash := crypto.Keccak256([]byte(string("totalSupply()")))
	data := "0x" + hex.EncodeToString(functionHash[0:4])
	req := request{usdt, data}

	var result string
	if err := client.Call(&result, "eth_call", req, "latest"); err != nil {
		log.Fatal(err)
	}
	totalSupply := big.NewInt(0)
	totalSupply, ok := totalSupply.SetString(result[2:], 16)
	if !ok {
		log.Fatal("err")
		return
	}
	fmt.Println("usdt total supply =", totalSupply)
}
