package node

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/martinmunillas/munichain/munichain"
)

type StatusRes struct {
	Hash                munichain.Hash          `json:"block_hash"`
	Number              uint64                  `json:"block_number"`
	KnownPeers          map[string]PeerNode     `json:"peers_known"`
	PendingTransactions []munichain.Transaction `json:"pending_transactions"`
}

func nodeStatusHandler(w http.ResponseWriter, r *http.Request, n *Node) {
	err := writeRes(w, StatusRes{
		Hash:                n.state.LatestBlockHash,
		Number:              n.state.LatestBlock.Header.Number,
		KnownPeers:          n.KnownPeers,
		PendingTransactions: n.getPendingTransactionsArray(),
	})
	if err != nil {
		writeErrRes(w, err)
	}

}

type BalancesRes struct {
	Hash   munichain.Hash     `json:"block_hash"`
	Amount munichain.Balances `json:"amount"`
}

func listBalancesHandler(w http.ResponseWriter, s *munichain.State) {
	err := writeRes(w, BalancesRes{s.LatestBlockHash, s.Balances})
	if err != nil {
		writeErrRes(w, err)
	}
}

type AddTransactionReq struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount uint   `json:"amount"`
}

type AddTransactionRes struct {
	Success bool `json:"success"`
}

func addTransactionHandler(w http.ResponseWriter, r *http.Request, n *Node) {
	var req []AddTransactionReq
	err := readReq(r, &req)
	if err != nil {
		writeErrRes(w, err)
		return
	}

	var txs []munichain.Transaction
	for _, tx := range req {
		txs = append(txs, munichain.Transaction{
			From:   munichain.Address(tx.From),
			To:     munichain.Address(tx.To),
			Amount: tx.Amount,
		})
	}

	for _, tx := range txs {
		n.addPendingTransaction(tx)
	}
	if err != nil {
		writeErrRes(w, err)
		return
	}

	writeRes(w, AddTransactionRes{true})
}

func syncHandler(w http.ResponseWriter, r *http.Request, dataDir string) {
	reqHash := r.URL.Query().Get(syncFromBlockQueryKey)
	hash := munichain.Hash{}
	err := hash.UnmarshalText([]byte(reqHash))
	if err != nil {
		writeErrRes(w, err)
		return
	}
	blocks, err := munichain.GetBlocksAfter(hash, dataDir)
	if err != nil {
		writeErrRes(w, err)
		return
	}
	writeRes(w, blocks)
}

type AddPeerRes struct {
	Success bool `json:"success"`
}

func joinPeerHandler(w http.ResponseWriter, r *http.Request, n *Node) {
	peerIP := r.URL.Query().Get(joinPeerIPQueryKey)
	peerPortRaw := r.URL.Query().Get(joinPeerPortQueryKey)

	peerPort, err := strconv.ParseUint(peerPortRaw, 10, 32)
	if err != nil {
		writeRes(w, AddPeerRes{false})
		return
	}

	peer := PeerNode{
		IP:          peerIP,
		Port:        peerPort,
		IsBootstrap: false,
		IsActive:    true,
	}

	if !n.isKnownPeer(peer) {
		n.addPeer(peer)

		fmt.Printf("Peer '%s' was added into KnownPeers\n", peer.TcpAddress())
	}

	writeRes(w, AddPeerRes{true})
}
