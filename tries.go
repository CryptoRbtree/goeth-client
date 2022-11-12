package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	ethtrie "github.com/ethereum/go-ethereum/trie"
)

func tries() {
	diskdb := memorydb.New()
	triedb := ethtrie.NewDatabase(diskdb)
	trie, err := ethtrie.New(common.Hash{}, common.Hash{}, triedb)
	if err != nil {
		log.Fatal(err)
	}

	trie.Update([]byte("foo"), []byte("bar"))
	trie.Update([]byte("eth"), []byte("pos"))
	value, err := trie.TryGet([]byte("foo"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("find match:", string(value))
}
