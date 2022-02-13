package munichain

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

type Address string

type Balances = map[Address]uint

type State struct {
	Balances        Balances
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

		err := state.applyBlock(blockFs.Value)
		if err != nil {
			return nil, err
		}

		state.LatestBlockHash = blockFs.Key
		state.LatestBlock = blockFs.Value
	}

	return state, nil
}

func (s *State) AddBlocks(blocks []Block) error {
	for _, b := range blocks {
		hash, err := s.AddBlock(b)
		if err != nil {
			return err
		}
		fmt.Printf("Added block %s\n", hash.ToString())
	}
	return nil
}

func (s *State) AddBlock(b Block) (Hash, error) {
	c := s.copy()

	err := c.applyBlock(b)
	if err != nil {
		return Hash{}, err
	}

	blockHash, err := b.Hash()
	if err != nil {
		return Hash{}, err
	}

	blockFs := BlockFS{
		Key:   blockHash,
		Value: b,
	}

	blockFsJson, err := json.Marshal(blockFs)
	if err != nil {
		return Hash{}, err
	}

	fmt.Printf("Persisting Block to disk:")
	fmt.Printf("\t%s\n", blockFsJson)

	_, err = s.dbFile.Write(append(blockFsJson, '\n'))
	if err != nil {
		return Hash{}, err
	}

	s.Balances = c.Balances
	s.LatestBlockHash = blockHash
	s.LatestBlock = b

	return blockHash, nil
}

func (s *State) copy() *State {
	c := &State{
		LatestBlockHash: s.LatestBlockHash,
		LatestBlock:     s.LatestBlock,
		Balances:        make(map[Address]uint),
	}
	for acc, balance := range s.Balances {
		c.Balances[acc] = balance
	}
	return c

}

func (s *State) applyBlock(b Block) error {
	isGenesis := s.LatestBlockHash == Hash{}
	nextExpectedBlockNumber := s.LatestBlock.Header.Number + 1

	if isGenesis && b.Header.Number != 0 && b.Header.Number != nextExpectedBlockNumber {
		return fmt.Errorf(
			"next expected block must be '%d' not '%d'",
			nextExpectedBlockNumber,
			b.Header.Number,
		)
	}
	// validate the incoming block parent hash equals
	// the current (latest known) hash
	if s.LatestBlock.Header.Number > 0 && !reflect.DeepEqual(b.Header.Previous, s.LatestBlockHash) {
		return fmt.Errorf(
			"next block previous hash must be '%x' not '%x'",
			s.LatestBlockHash,
			b.Header.Previous,
		)
	}
	hash, err := b.Hash()
	if err != nil {
		return err
	}
	if !IsBlockHashValid(hash) {
		return fmt.Errorf("invalid block hash %x", hash)
	}

	err = s.applyTransactions(b.Transactions)
	if err != nil {
		return err
	}

	s.Balances[b.Header.Miner] += BlockReward

	return nil
}

func (s *State) applyTransactions(transactions []Transaction) error {
	for _, tx := range transactions {
		if err := s.applyTransaction(tx); err != nil {
			return err
		}
	}
	return nil
}

func (s *State) applyTransaction(tx Transaction) error {
	if tx.isPrinting() {
		s.Balances[tx.To] += tx.Amount
		return nil
	}

	if !tx.isValid() {
		return fmt.Errorf("invalid transaction")
	}

	if tx.Amount > s.Balances[tx.From] {
		return fmt.Errorf("wrong transaction. '%s' balance is %d, but trying to send %d", tx.From, s.Balances[tx.From], tx.Amount)
	}

	s.Balances[tx.From] -= tx.Amount
	s.Balances[tx.To] += tx.Amount

	return nil
}
func (state *State) Close() {
	state.dbFile.Close()
}
