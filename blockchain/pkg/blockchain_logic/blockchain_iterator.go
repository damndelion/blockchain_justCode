package blockchainlogic

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
)

type BlockchainIterator struct {
	currentHash string
	db          *sql.DB
}

func (i *BlockchainIterator) Next() *Block {
	var block Block
	var transactionsJSON string

	// Retrieve the block data from the 'blocks' table using the current hash
	row := i.db.QueryRow("SELECT hash, transactions, previous_hash, timestamp, nonce FROM blocks WHERE hash = $1", i.currentHash)

	err := row.Scan(&block.Hash, &transactionsJSON, &block.PrevHash, &block.Timestamp, &block.Nonce)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil // End of the blockchain
		}

		log.Panic(err)
	}

	if err = json.Unmarshal([]byte(transactionsJSON), &block.Transactions); err != nil {
		log.Fatal(err)

		return nil
	}

	i.currentHash = block.PrevHash

	return &block
}
