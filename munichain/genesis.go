package munichain

var genesisBlock = Block{
	Header: BlockHeader{
		Time:     1642358385,
		Previous: Hash{},
	},
	Transactions: []Transaction{
		NewTransaction("munichain", "martinmunilla", 100_000_000),
	},
}

func (genesis *Block) isPrintingTx(tx Transaction) bool {
	return tx.From == "munichain"
}
