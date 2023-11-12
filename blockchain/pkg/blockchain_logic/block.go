package blockchain_logic

import (
	"bytes"
	"crypto/sha256"
	"time"
)

type Block struct {
	Hash         string
	Transactions []*Transaction
	PrevHash     string
	Timestamp    time.Time
	Nonce        int
}

func CreateBlock(transactions []*Transaction, prevHash string) *Block {
	block := &Block{"", transactions, prevHash, time.Now(), 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func Genesis(coinbase *Transaction) *Block {

	return CreateBlock([]*Transaction{coinbase}, "0")
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}
