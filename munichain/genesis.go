package munichain

var GenesisBlock = Block{
	Header: BlockHeader{
		Time:     1642358385,
		Number:   0,
		Previous: Hash{},
		Nonce:    77736833,
		Miner:    "martinmunilla",
	},
	Transactions: []Transaction{
		{
			From:   "munichain",
			To:     "martinmunilla",
			Amount: 100_000_000,
			Data:   "genesis",
			Time:   1642358385,
		},
	},
}

func (tx *Transaction) isPrinting() bool {
	return tx.From == "munichain"
}
