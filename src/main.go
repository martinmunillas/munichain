package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Address string

type Transaction struct {
	from    Address `json:"from"`
	to      Address `json:"to"`
	amount  uint    `json:"amount"`
	message string  `json:"message"`
}

type State struct {
	balances map[Address]uint
	memPool  []Transaction

	transactionsFile []byte
}

func (state *State) apply(tx *Transaction) {

	if state.isValidTransaction(tx) {
		return
	}

	state.balances[tx.from] -= tx.amount
	state.balances[tx.to] += tx.amount
}

func (state *State) isValidTransaction(tx *Transaction) bool {
	if tx.amount <= 0 {
		return false

	if state.balances[tx.from] < tx.amount {
		return false
	}

	return true
}

func newStateFromDisk() (*State, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	genesisFilePath := filepath.Join(cwd, "db", "genesis.json")

	genesis, err := loadGenesis(genesisFilePath)
	if err != nil {
		return nil, err
	}

	balances := make(map[Address]uint)
	for address, balance := range genesis.Balances {
		balances[address] = balance
	}

	transactionsFilePath := filepath.Join(cwd, "db", "transactions.json")

	transactionsFile, err := ioutil.ReadFile(transactionsFilePath)
	if err != nil {
		return nil, err
	}

	state := &State{
		balances:         balances,
		memPool:          make([]Transaction, 0),
		transactionsFile: transactionsFile,
	}

	var transactions []Transaction
	json.Unmarshal(transactionsFile, &transactions)

	for transaction := range transactions {
		state.apply(transaction)
	}

	return state, nil
}
