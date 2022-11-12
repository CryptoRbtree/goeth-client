package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	store "github.com/CryptoRbtree/goeth-client/contract-store"
)

func contractReadWrite() {
	var (
		ctx         = context.Background()
		url         = "https://eth-goerli.g.alchemy.com/v2/" + os.Getenv("ALCHEMY_ID")
		client, err = ethclient.DialContext(ctx, url)
	)
	if err != nil {
		log.Fatal(err)
	}

	// 1. load contract
	address := common.HexToAddress("0xd3047d5bbcbcfe4256f9c1668c57b3d875c4adea")
	instance, err := store.NewStore(address, client)
	if err != nil {
		log.Fatal(err)
	}

	// 2. read contract
	version, err := instance.Version(nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("store contract version =", version)

	// 3. write contract
	privateKey, err := crypto.HexToECDSA(os.Getenv("ACCOUNT_KEY"))
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	gasTipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	//auth.GasPrice = big.NewInt(1e9) // only for legacy transactions
	auth.GasFeeCap = gasPrice
	auth.GasTipCap = gasTipCap

	key := [32]byte{}
	value := [32]byte{}
	copy(key[:], []byte("foo"))
	copy(value[:], []byte("bar"))
	if err != nil {
		log.Fatal(err)
	}

	// estimate gas
	parsed, err := abi.JSON(strings.NewReader(store.StoreABI))
	if err != nil {
		log.Fatal(err)
	}
	encodedData, err := parsed.Pack("setItem", key, value)
	if err != nil {
		log.Fatal(err)
	}
	estimatedGas, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:      fromAddress,
		To:        &address,
		Data:      encodedData,
		GasFeeCap: gasPrice,
		GasTipCap: gasTipCap,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	auth.GasLimit = estimatedGas // in units

	tx, err := instance.SetItem(auth, key, value)

	txHash := common.HexToHash(tx.Hash().Hex())
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	fmt.Println("contract address =", address.Hex())
	fmt.Println("transaction hash:", tx.Hash().Hex())
	fmt.Println("transaction gas limit:", tx.Gas())
	fmt.Println("transaction fee cap per gas:", tx.GasFeeCap())
	fmt.Println("transaction tip cap per gas:", tx.GasTipCap())
	fmt.Println("transaction data:", hex.EncodeToString(tx.Data()))
	fmt.Println("transaction to:", tx.To().Hex())
	fmt.Println("isPending:", isPending)

	for isPending {
		time.Sleep(time.Second * 5)
		fmt.Println("pending...")
		_, isPending, _ = client.TransactionByHash(context.Background(), txHash)
	}

	result, err := instance.Items(nil, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(result[:])) // "bar"
}
