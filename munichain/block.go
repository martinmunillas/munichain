package munichain

import (
	"crypto/sha256"
	"encoding/json"
)

type Hash [32]byte

type Block struct {
	Header       BlockHeader
	Transactions []Transaction
}

type BlockHeader struct {
	Previous Hash
	Time     uint64
}

func (block *Block) Hash() (Hash, error) {
	data, err := json.Marshal(block)
	if err != nil {
		return Hash{}, err
	}
	return sha256.Sum256(data), nil
}
