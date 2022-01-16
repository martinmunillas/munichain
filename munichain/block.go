package munichain

import (
	"crypto/sha256"
	"encoding/json"
)

type Hash [32]byte

type Block struct {
	Header       BlockHeader   `json:"header"`
	Transactions []Transaction `json:"transactions"`
}

type BlockHeader struct {
	Previous Hash   `json:"previous"`
	Time     uint64 `json:"time"`
}

type BlockFS struct {
	Key   Hash  `json:"hash"`
	Value Block `json:"block"`
}

func (block *Block) Hash() (Hash, error) {
	data, err := json.Marshal(block)
	if err != nil {
		return Hash{}, err
	}
	return sha256.Sum256(data), nil
}

func (block *Block) getBalances() Balances {
	balances := Balances{}
	for _, tx := range block.Transactions {
		if block.isGenesisTx(tx) {
			balances[tx.To] += tx.Amount
		} else {
			balances[tx.From] -= tx.Amount
			balances[tx.To] += tx.Amount
		}
	}
	return balances
}
