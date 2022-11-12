package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	store "github.com/CryptoRbtree/goeth-client/contract-store"
)

func contractDeploy() {
	var (
		ctx         = context.Background()
		url         = "https://eth-goerli.g.alchemy.com/v2/" + os.Getenv("ALCHEMY_ID")
		client, err = ethclient.DialContext(ctx, url)
	)
	if err != nil {
		log.Fatal(err)
	}

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
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(500000) // in units
	//auth.GasPrice = big.NewInt(1e9) // only for legacy transactions
	auth.GasFeeCap = gasPrice
	auth.GasTipCap = gasTipCap

	input := "1.0"
	address, tx, instance, err := store.DeployStore(auth, client, input)
	if err != nil {
		log.Fatal(err)
	}

	txHash := common.HexToHash(tx.Hash().Hex())
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	fmt.Println("contract address =", address.Hex())
	fmt.Println("transaction hash:", tx.Hash().Hex())
	fmt.Println("transaction gas limit:", tx.Gas())
	fmt.Println("transaction fee cap per gas:", tx.GasFeeCap())
	fmt.Println("transaction tip cap per gas:", tx.GasTipCap())
	fmt.Println("transaction data:", hex.EncodeToString(tx.Data()))
	fmt.Println("transaction to:", tx.To())
	fmt.Println("isPending:", isPending)

	_ = instance
}
