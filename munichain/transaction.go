package munichain

import (
	"crypto/sha256"
	"encoding/json"
)

type Transaction struct {
	From   Address `json:"from"`
	To     Address `json:"to"`
	Amount uint    `json:"amount"`
	Data   string  `json:"data"`
	Time   uint64  `json:"time"`
}

func (tx *Transaction) isValid() bool {
	return tx.Amount > 0
}

func (tx *Transaction) Hash() (Hash, error) {
	data, err := json.Marshal(tx)
	if err != nil {
		return Hash{}, err
	}
	return sha256.Sum256(data), nil
}
