package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

func rlpEncode() {
	var (
		ctx         = context.Background()
		url         = "https://eth-mainnet.g.alchemy.com/v2/" + os.Getenv("ALCHEMY_ID")
		client, err = ethclient.DialContext(ctx, url)
	)
	if err != nil {
		log.Fatal(err)
	}

	txHash := common.HexToHash("0xec9db5bfbcd30ad2e3070b626ed4f78abce88687c5d1eb23464242be5edcb537")
	tx, _, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("----original transaction----")
	fmt.Println("transaction hash:", tx.Hash().Hex())
	fmt.Println("transaction value:", tx.Value().String())
	fmt.Println("transaction gas limit:", tx.Gas())
	fmt.Println("transaction fee cap per gas:", tx.GasFeeCap())
	fmt.Println("transaction tip cap per gas:", tx.GasTipCap())
	fmt.Println("transaction gas price:", tx.GasPrice())
	fmt.Println("transaction nonce:", tx.Nonce())
	fmt.Println("transaction data:", hex.EncodeToString(tx.Data()))
	fmt.Println("transaction to:", tx.To().Hex())

	v, r, s := tx.RawSignatureValues()
	R := r.Bytes()
	S := s.Bytes()
	V := byte(v.Uint64())
	sig := make([]byte, 65)
	copy(sig[32-len(R):32], R)
	copy(sig[64-len(S):64], S)
	sig[64] = V
	fmt.Println(hex.EncodeToString(sig))

	rawTx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("rawTx = %x\n", rawTx)

	tx = new(types.Transaction)
	rawTxBytes, err := hex.DecodeString(hex.EncodeToString(rawTx))
	rlp.DecodeBytes(rawTxBytes, &tx)
	fmt.Println("----decode transaction----")
	fmt.Println("transaction hash:", tx.Hash().Hex())
	fmt.Println("transaction value:", tx.Value().String())
	fmt.Println("transaction gas limit:", tx.Gas())
	fmt.Println("transaction fee cap per gas:", tx.GasFeeCap())
	fmt.Println("transaction tip cap per gas:", tx.GasTipCap())
	fmt.Println("transaction gas price:", tx.GasPrice())
	fmt.Println("transaction nonce:", tx.Nonce())
	fmt.Println("transaction data:", hex.EncodeToString(tx.Data()))
	fmt.Println("transaction to:", tx.To().Hex())
}
