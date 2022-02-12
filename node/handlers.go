package node

import (
	"net/http"

	"github.com/martinmunillas/munichain/munichain"
)

type BalancesRes struct {
	Hash   munichain.Hash     `json:"block_hash"`
	Amount munichain.Balances `json:"amount"`
}

type AddTransactionReq struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount uint   `json:"amount"`
}

func nodeStatusHandler(w http.ResponseWriter, r *http.Request, n *Node) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\":\"ok\"}"))
}

func listBalancesHandler(w http.ResponseWriter, s *munichain.State) {
	err := writeRes(w, BalancesRes{s.LatestBlockHash, s.Balances})
	if err != nil {
		writeErrRes(w, err)
	}
}

func addTransactionHandler(w http.ResponseWriter, r *http.Request, s *munichain.State) {
	var req []AddTransactionReq
	err := readReq(r, &req)
	if err != nil {
		writeErrRes(w, err)
		return
	}

	var txs []*munichain.Transaction
	for _, req := range req {
		txs = append(txs, &munichain.Transaction{
			From:   munichain.Address(req.From),
			To:     munichain.Address(req.To),
			Amount: req.Amount,
		})
	}

	err = s.AddTransactions(txs...)
	if err != nil {
		writeErrRes(w, err)
		return
	}
	hash, err := s.Persist()
	if err != nil {
		writeErrRes(w, err)
		return
	}

	writeRes(w, hash)
}
