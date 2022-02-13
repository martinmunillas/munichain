package munichain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Hash [32]byte

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}

func (h Hash) ToString() string {
	return hex.EncodeToString(h[:])
}

func (h *Hash) UnmarshalText(data []byte) error {
	_, err := hex.Decode(h[:], data)
	return err
}

type Block struct {
	Header       BlockHeader   `json:"header"`
	Transactions []Transaction `json:"transactions"`
}

type BlockHeader struct {
	Previous Hash   `json:"previous"`
	Number   uint64 `json:"number"`
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
