package node

import (
	"fmt"
	"net/http"

	"github.com/martinmunillas/munichain/munichain"
)

const DefaultHttpPort = 8080

type Node struct {
	DataDir string
	Port    uint64

	s *munichain.State

	KnownPeers []PeerNode
}

type PeerNode struct {
	IP          string `json:"ip"`
	Port        uint64 `json:"port"`
	IsBootstrap bool   `json:"isBootstrap"`
	IsActive    bool   `json:"isActive"`
}

func (n *Node) Run() error {
	fmt.Printf("Listening on port %d\n", n.Port)
	state, err := munichain.NewStateFromDisk(n.DataDir)
	if err != nil {
		return err
	}
	defer state.Close()

	n.s = state

	http.HandleFunc("/balances/list", func(w http.ResponseWriter, r *http.Request) {
		listBalancesHandler(w, state)
	})
	http.HandleFunc("/transactions/add", func(w http.ResponseWriter, r *http.Request) {
		addTransactionHandler(w, r, state)
	})

	http.HandleFunc("/node/status", func(w http.ResponseWriter, r *http.Request) {
		nodeStatusHandler(w, r, n)
	})

	http.ListenAndServe(fmt.Sprintf(":%d", n.Port), nil)
	return nil
}
