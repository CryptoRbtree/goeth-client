package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

func transactionSignature() {
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

	v, r, s := tx.RawSignatureValues()
	R := r.Bytes()
	S := s.Bytes()
	V := byte(v.Uint64())
	sig := make([]byte, 65)
	copy(sig[32-len(R):32], R)
	copy(sig[64-len(S):64], S)
	sig[64] = V

	signer := types.NewLondonSigner(tx.ChainId())
	hash := signer.Hash(tx)

	rawTx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", rawTx)

	publicKeyBytes, err := crypto.Ecrecover(hash.Bytes(), sig)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("publicKey = %x\n", publicKeyBytes)

	publicKeyECDSA, err := crypto.UnmarshalPubkey(publicKeyBytes)
	if err != nil {
		log.Fatal(err)
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("address =", address) // 0x5827c0ccf705720cfa395e3fb2dcc449aeef331c
}
