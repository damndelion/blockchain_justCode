package blockchain_logic

import (
	"bytes"
	"crypto/rsa"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type Blockchain struct {
	tip string
	Db  *sql.DB
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
	} else {
		return NewBlockchain(db, address)

	}

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
	var lastHash string

	// Retrieve the hash of the latest block
	err := bc.Db.QueryRow("SELECT hash FROM blocks ORDER BY id DESC LIMIT 1").Scan(&lastHash)
	if err != nil {
		return err
	}
	newBlock := CreateBlock(transactions, lastHash)
	transactionsJSON, err := json.Marshal(transactions)
	// Insert the new block into the database
	_, err = bc.Db.Exec("INSERT INTO blocks (hash, transactions, previous_hash, timestamp, nonce) VALUES ($1, $2, $3, $4, $5)",
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

			if tx.IsCoinbase() == false {
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

// Iterator returns a BlockchainIterat
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.Db}

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

	return Transaction{}, errors.New("transaction is not found")
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

	chain := NewBlockchain(bc.Db, address)

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
	if !ValidateAddress(from) {
		return fmt.Errorf("ERROR: Sender address is not valid")
	}
	if !ValidateAddress(to) {
		return fmt.Errorf("ERROR: Sender address is not valid")
	}

	NewBlockchain(bc.Db, from)

	tx, err := NewUTXOTransaction(from, to, amount, bc)
	if err != nil {
		return err
	}
	err = bc.MineBlock([]*Transaction{tx})
	if err != nil {
		return err
	}

	return nil

}

type CoinGeckoResponse struct {
	Bitcoin struct {
		USD float64 `json:"usd"`
	} `json:"bitcoin"`
}

func (bc *Blockchain) GetBalanceInUSD(address string) (float64, error) {
	url := "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd"

	response, err := http.Get(url)
	if err != nil {
		return -1, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)

	var data CoinGeckoResponse

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return -1, err
	}

	bitcoinPriceUSD := data.Bitcoin.USD

	bitcoinBalance := bc.GetBalance(address)
	totalBalanceUSD := bitcoinBalance * bitcoinPriceUSD

	return totalBalanceUSD, nil
}
