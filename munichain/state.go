package munichain

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Address string

type Balances = map[Address]uint

type State struct {
	Balances        Balances
	memPool         []Transaction
	LatestBlockHash Hash

	dbFile *os.File
}

func NewStateFromDisk(dataDir string) (*State, error) {
	err := initDataDirIfNotExists(dataDir)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(getBlocksFilePath(dataDir), os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	hash, err := genesisBlock.Hash()
	if err != nil {
		return nil, err
	}

	state := &State{
		Balances:        genesisBlock.getBalances(),
		LatestBlockHash: hash,
		memPool:         []Transaction{},
		dbFile:          file,
	}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		blockFsJson := scanner.Text()
		blockFs := BlockFS{}
		err = json.Unmarshal([]byte(blockFsJson), &blockFs)
		if err != nil {
			return nil, err
		}

		state.applyBlock(blockFs.Value)

		state.LatestBlockHash = blockFs.Key

	}

	return state, nil
}

func (state *State) AddTransactions(txs ...*Transaction) {
	for _, tx := range txs {
		state.memPool = append(state.memPool, *tx)
		state.applyTransaction(tx)
	}
}

func (state State) applyBlock(block Block) {
	for _, tx := range block.Transactions {
		state.applyTransaction(&tx)
	}
}

func (state *State) applyTransaction(tx *Transaction) error {

	if state.isValidTransaction(tx) {
		return fmt.Errorf("invalid transaction: %v", tx)
	}

	state.Balances[tx.From] -= tx.Amount
	state.Balances[tx.To] += tx.Amount
	return nil
}

func (state *State) isValidTransaction(tx *Transaction) bool {
	if tx.Amount <= 0 {
		return false
	}

	if state.Balances[tx.From] < tx.Amount {
		return false
	}

	return true
}

func (state *State) Persist() (Hash, error) {
	block := &Block{
		Header: BlockHeader{
			Previous: state.LatestBlockHash,
			Time:     uint64(time.Now().Unix()),
		},
		Transactions: state.memPool,
	}

	blockHash, err := block.Hash()
	if err != nil {
		return Hash{}, err
	}

	blockFs := BlockFS{
		Key:   blockHash,
		Value: *block,
	}

	blockFsJson, err := json.Marshal(blockFs)
	if err != nil {
		return Hash{}, err
	}

	fmt.Printf("Persisting new Block to disk:\n")
	fmt.Printf("\t%s\n", blockFsJson)

	_, err = state.dbFile.Write(append(blockFsJson, '\n'))
	if err != nil {
		return Hash{}, err
	}
	state.LatestBlockHash = blockHash
	state.memPool = []Transaction{}
	return blockHash, nil

}
