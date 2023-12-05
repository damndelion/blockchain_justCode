package blockchainlogic

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type Blockchain struct {
	tip string
	DB  *sql.DB
	mu  *sync.Mutex
}

func CreateBlockchain(db *sql.DB, address string) *Blockchain {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	// Check if the blocks table is empty
	var blockCount int
	err := db.QueryRow("SELECT COUNT(*) FROM blocks").Scan(&blockCount)
	if err != nil {
		log.Fatal(err)
	}
	if blockCount == 0 {
		var lastHash string

		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		genesis := Genesis(cbtx)
		jsonData, err := json.Marshal(genesis.Transactions)
		if err != nil {
			log.Fatal(err)
		}
		transactionsString := string(jsonData)

		_, err = db.Exec("INSERT INTO blocks (hash, transactions, previous_hash, timestamp, nonce) VALUES ($1, $2, $3, $4, $5)",
			genesis.Hash, transactionsString, genesis.PrevHash, genesis.Timestamp, 0)
		if err != nil {
			log.Fatal(err)
		}

		lastHash = genesis.Hash
		mu := &sync.Mutex{}
		chain := Blockchain{lastHash, db, mu}

		return &chain
	}

	return NewBlockchain(db, address)
}

func NewBlockchain(db *sql.DB, address string) *Blockchain {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	var lastHash string
	var lastBlockHash string
	err := db.QueryRow("SELECT hash FROM blocks ORDER BY id DESC LIMIT 1").Scan(&lastBlockHash)
	if err != nil {
		log.Fatal(err)
	}
	lastHash = lastBlockHash
	if err != nil {
		log.Panic(err)
	}
	mu := &sync.Mutex{}
	chain := Blockchain{lastHash, db, mu}

	return &chain
}

func (bc *Blockchain) MineBlock(transactions []*Transaction) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	var lastHash string
	// Retrieve the hash of the latest block
	err := bc.DB.QueryRow("SELECT hash FROM blocks ORDER BY id DESC LIMIT 1").Scan(&lastHash)
	if err != nil {
		return err
	}
	newBlock := CreateBlock(transactions, lastHash)
	transactionsJSON, err := json.Marshal(transactions)
	if err != nil {
		return err
	}
	// Insert the new block into the database
	_, err = bc.DB.Exec("INSERT INTO blocks (hash, transactions, previous_hash, timestamp, nonce) VALUES ($1, $2, $3, $4, $5)",
		newBlock.Hash, string(transactionsJSON), newBlock.PrevHash, newBlock.Timestamp, newBlock.Nonce)

	if err != nil {
		log.Fatal(err)
	}

	bc.tip = newBlock.Hash

	return nil
}

func (bc *Blockchain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]float64)
	bci := bc.Iterator()

	for {
		block := bci.Next()
		if block == nil {
			break
		}
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == float64(outIdx) {
							continue Outputs
						}
					}
				}

				if out.IsLockedWithKey(pubKeyHash) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			if !tx.IsCoinbase() {
				for _, in := range tx.Vin {
					if in.UsesKey(pubKeyHash) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return unspentTXs
}

// Iterator returns a BlockchainIterator.
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.DB}

	return bci
}

func (bc *Blockchain) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnspentTransactions(pubKeyHash)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (bc *Blockchain) FindSpendableOutputs(pubKeyHash []byte, amount float64) (float64, map[string][]float64) {
	unspentOutputs := make(map[string][]float64)
	unspentTXs := bc.FindUnspentTransactions(pubKeyHash)
	accumulated := 0.0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], float64(outIdx))

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New(fmt.Sprintf("transaction is not found"))
}

func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}

func (bc *Blockchain) SignTransaction(tx *Transaction, privKey *rsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

func (bc *Blockchain) GetBalance(address string) float64 {
	chain := NewBlockchain(bc.DB, address)

	balance := 0.0
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := chain.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	return balance
}

func (bc *Blockchain) Send(from, to string, amount float64) error {
	key := generateTransactionKey(from, to, amount)
	if !ValidateAddress(from) {
		return errors.New(fmt.Sprintf("ERROR: Sender address is not valid"))
	}
	if !ValidateAddress(to) {
		return errors.New(fmt.Sprintf("ERROR: Sender address is not valid"))
	}

	NewBlockchain(bc.DB, from)

	tx, err := NewUTXOTransaction(from, to, amount, bc, key)
	if err != nil {
		return err
	}
	err = bc.MineBlock([]*Transaction{tx})
	if err != nil {
		return err
	}

	return nil
}

func generateTransactionKey(from, to string, amount float64) string {
	amountStr := fmt.Sprintf("%.2f", amount)

	timestamp := time.Now().Unix() / 60

	data := from + to + amountStr + strconv.FormatInt(timestamp, 10)

	hash := sha256.New()
	hash.Write([]byte(data))
	hashCode := fmt.Sprintf("%x", hash.Sum(nil))

	return hashCode
}
