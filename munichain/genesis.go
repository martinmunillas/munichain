package munichain

var genesisBlock = Block{
	Header: BlockHeader{
		Time:     1642358385,
		Number:   0,
		Previous: Hash{},
	},
	Transactions: []Transaction{
		NewTransaction("munichain", "martinmunilla", 100_000_000),
	},
}

func (tx *Transaction) isPrinting() bool {
	return tx.From == "munichain"
}
