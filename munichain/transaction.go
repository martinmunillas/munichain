package munichain

type Transaction struct {
	From     Address `json:"from"`
	To       Address `json:"to"`
	Amount   uint    `json:"amount"`
	Message  string  `json:"message"`
	Rejected bool    `json:"rejected"`
}
