package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/martinmunillas/munichain/munichain"
)

const httpPort = 8080

type BalancesRes struct {
	Hash   munichain.Hash     `json:"block_hash"`
	Amount munichain.Balances `json:"amount"`
}

type AddTransactionReq struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount uint   `json:"amount"`
}

func Run(dataDir string) error {
	state, err := munichain.NewStateFromDisk(dataDir)
	if err != nil {
		return err
	}
	defer state.Close()
	http.HandleFunc("/balances/list", func(w http.ResponseWriter, r *http.Request) {
		listBalancesHandler(w, state)
	})
	http.HandleFunc("/transactions/add", func(w http.ResponseWriter, r *http.Request) {
		addTransactionHandler(w, r, state)
	})
	http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
	return nil
}

func listBalancesHandler(w http.ResponseWriter, state *munichain.State) {
	err := writeRes(w, BalancesRes{state.LatestBlockHash, state.Balances})
	if err != nil {
		writeErrRes(w, err)
	}
}

func addTransactionHandler(w http.ResponseWriter, r *http.Request, state *munichain.State) {
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

	err = state.AddTransactions(txs...)
	if err != nil {
		writeErrRes(w, err)
		return
	}
	hash, err := state.Persist()
	if err != nil {
		writeErrRes(w, err)
		return
	}

	writeRes(w, hash)
}

func writeRes(w http.ResponseWriter, res interface{}) error {
	data, err := json.Marshal(res)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return nil
}

func writeErrRes(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func readReq(r *http.Request, reqBody interface{}) error {
	reqBodyJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("unable to read request body. %s", err.Error())
	}
	defer r.Body.Close()

	err = json.Unmarshal(reqBodyJson, reqBody)
	if err != nil {
		return fmt.Errorf("unable to unmarshal request body. %s", err.Error())
	}

	return nil
}
