package munichain

type Address string

type Balances = map[Address]uint

type Meta struct {
	LastSync uint64
}
type State struct {
	Balances     Balances
	Transactions []Transaction
	Meta         Meta
}

func (state *State) sync() {
	for _, tx := range state.Transactions[state.Meta.LastSync:] {
		state.apply(&tx)
	}
	state.Meta.LastSync = uint64(len(state.Transactions) - 1)
}

func (state *State) Add(tx *Transaction) {
	state.Transactions = append(state.Transactions, *tx)
	state.sync()
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

func (state *State) Persist() error {
	if err := writeJson(state.Balances, "db", "balances.json"); err != nil {
		return err
	}
	return nil
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
		Balances:     genesis,
		Transactions: transactions,
		Meta: Meta{
			LastSync: 0,
		},
	}

	state.sync()

	return state, nil
}
