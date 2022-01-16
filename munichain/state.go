package munichain

import "time"

type Address string

type Balances = map[Address]uint

type State struct {
	Balances    Balances
	memPool     []Transaction
	currentHash Hash
}

func (state *State) AddTransactions(txs ...*Transaction) {
	for tx := range txs {
	}
}

func (state *State) apply(tx *Transaction) {

	if state.isValidTransaction(tx) {
		tx.Rejected = true
		return
	}

	state.Balances[tx.From] -= tx.Amount
	state.Balances[tx.To] += tx.Amount
}

func (state *State) isValidTransaction(tx *Transaction) bool {
	if tx.Amount <= 0 {
		return false
	}

	if state.Balances[tx.From] < tx.Amount {
		return false
	}

	return true
}

func (state *State) Persist() (Hash, error) {
	block := &Block{
		Header: BlockHeader{
			Previous: state.currentHash,
			Time:     uint64(time.Now().Unix()),
		},
		Transactions: state.memPool,
	}

	blockHash, err := block.Hash()
	if err != nil {
		return Hash{}, err
	}

}

func NewStateFromDisk() (*State, error) {

	var genesis Balances

	err := loadJson(&genesis, "genesis.json")
	if err != nil {
		return nil, err
	}

	var transactions []Transaction

	err = loadJson(&transactions, "db", "transactions.json")
	if err != nil {
		return nil, err
	}
	state := &State{
		Balances: genesis,
	}

	return state, nil
}
