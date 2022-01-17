package munichain

type Transaction struct {
	From   Address `json:"from"`
	To     Address `json:"to"`
	Amount uint    `json:"amount"`
}

func NewTransaction(from Address, to Address, value uint) Transaction {
	return Transaction{
		From:   from,
		To:     to,
		Amount: value,
	}
}
func (tx *Transaction) isValid() bool {
	return tx.Amount > 0
}
