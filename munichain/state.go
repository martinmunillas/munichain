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
	LatestBlock     Block

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

	state := &State{
		Balances:        Balances{},
		LatestBlockHash: Hash{},
		LatestBlock:     Block{},
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

		var transactions []*Transaction
		for _, txFs := range blockFs.Value.Transactions {
			transactions = append(transactions, &txFs)
		}
		state.AddTransactions(transactions...)

		state.LatestBlockHash = blockFs.Key
		state.LatestBlock = blockFs.Value
	}

	return state, nil
}

func (state *State) AddTransactions(txs ...*Transaction) error {
	prevState := *state
	for _, tx := range txs {
		err := state.addTransaction(tx)
		if err != nil {
			state.restore(prevState)
			return err
		}
	}
	return nil
}

func (state *State) restore(old State) {
	state = &old
}

func (state *State) addTransaction(tx *Transaction) error {
	if genesisBlock.isPrintingTx(*tx) {
		state.Balances[tx.To] = tx.Amount
		return nil
	}
	if !tx.isValid() {
		return fmt.Errorf("invalid transaction: %v", tx)
	}
	if state.Balances[tx.From] < tx.Amount {
		return fmt.Errorf("insufficient funds: %v", tx)
	}

	state.Balances[tx.From] -= tx.Amount
	state.Balances[tx.To] += tx.Amount
	state.memPool = append(state.memPool, *tx)
	return nil
}

func (state *State) Persist() (Hash, error) {
	block := &Block{
		Header: BlockHeader{
			Previous: state.LatestBlockHash,
			Number:   state.LatestBlock.Header.Number + 1,
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

func (state *State) Close() {
	state.dbFile.Close()
}
