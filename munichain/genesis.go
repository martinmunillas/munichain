package munichain

import "time"

var genesisBlock = Block{
	Header: BlockHeader{
		Time:     uint64(time.Now().Unix()),
		Previous: Hash{},
	},
	Transactions: []Transaction{
		NewTransaction("genesis", "martinmunilla", 100_000_000),
	},
}

func (genesis *Block) isGenesisTx(tx Transaction) bool {
	return tx.From == "genesis"
}
