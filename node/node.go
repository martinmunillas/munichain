package node

import (
	"context"
	"fmt"
	"net/http"

	"github.com/martinmunillas/munichain/munichain"
)

const DefaultHttpPort = 8080

type Node struct {
	DataDir string
	Port    uint64

	state *munichain.State

	KnownPeers map[string]PeerNode
}

func New(dataDir string, port uint64, bootstrap PeerNode) *Node {
	knownPeers := map[string]PeerNode{
		bootstrap.TcpAddress(): bootstrap,
	}
	return &Node{
		DataDir:    dataDir,
		Port:       port,
		KnownPeers: knownPeers,
	}
}

func (n *Node) Run() error {
	ctx := context.Background()
	fmt.Printf("Listening on port %d\n", n.Port)
	state, err := munichain.NewStateFromDisk(n.DataDir)
	if err != nil {
		return err
	}
	defer state.Close()

	n.state = state

	go n.sync(ctx)

	http.HandleFunc(listBalancesEndpoint, func(w http.ResponseWriter, r *http.Request) {
		listBalancesHandler(w, state)
	})
	http.HandleFunc(addTransactionsEndpoint, func(w http.ResponseWriter, r *http.Request) {
		addTransactionHandler(w, r, state)
	})

	http.HandleFunc(statusEndpoint, func(w http.ResponseWriter, r *http.Request) {
		nodeStatusHandler(w, r, n)
	})

	http.HandleFunc(syncEndpoint, func(w http.ResponseWriter, r *http.Request) {
		syncHandler(w, r, n.DataDir)
	})

	http.ListenAndServe(fmt.Sprintf(":%d", n.Port), nil)
	return nil
}
