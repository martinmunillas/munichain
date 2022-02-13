package node

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/martinmunillas/munichain/munichain"
)

const DefaultHttpPort = 8080

type Node struct {
	DataDir string
	Port    uint64
	IP      string
	Miner   munichain.Address

	state *munichain.State

	pendingTransactions  map[munichain.Hash]munichain.Transaction
	archivedTransactions map[munichain.Hash]munichain.Transaction
	isMining             bool
	newSyncedBlocks      chan munichain.Block

	KnownPeers map[string]PeerNode
}

func New(dataDir string, port uint64, bootstrap PeerNode, ip string, miner munichain.Address) *Node {
	knownPeers := map[string]PeerNode{
		bootstrap.TcpAddress(): bootstrap,
	}
	return &Node{
		DataDir: dataDir,
		Port:    port,
		IP:      ip,
		Miner:   miner,

		pendingTransactions:  make(map[munichain.Hash]munichain.Transaction),
		archivedTransactions: make(map[munichain.Hash]munichain.Transaction),
		isMining:             false,
		newSyncedBlocks:      make(chan munichain.Block),

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
	go n.mine(ctx)

	http.HandleFunc(listBalancesEndpoint, func(w http.ResponseWriter, r *http.Request) {
		listBalancesHandler(w, state)
	})
	http.HandleFunc(addTransactionsEndpoint, func(w http.ResponseWriter, r *http.Request) {
		addTransactionHandler(w, r, n)
	})

	http.HandleFunc(statusEndpoint, func(w http.ResponseWriter, r *http.Request) {
		nodeStatusHandler(w, r, n)
	})

	http.HandleFunc(syncEndpoint, func(w http.ResponseWriter, r *http.Request) {
		syncHandler(w, r, n.DataDir)
	})

	http.HandleFunc(joinPeerEndpoint, func(w http.ResponseWriter, r *http.Request) {
		joinPeerHandler(w, r, n)
	})

	http.ListenAndServe(fmt.Sprintf(":%d", n.Port), nil)
	return nil
}

func (n *Node) mine(ctx context.Context) error {
	var miningCtx context.Context
	var stopCurrentMining context.CancelFunc

	ticker := time.NewTicker(time.Second * 10)

	for {
		select {
		case <-ticker.C:
			go func() {
				if len(n.pendingTransactions) > 0 && !n.isMining {
					n.isMining = true

					miningCtx, stopCurrentMining = context.WithCancel(ctx)
					err := n.minePendingTransactions(miningCtx)
					if err != nil {
						fmt.Printf("ERROR: %s\n", err)
					}

					n.isMining = false
				}
			}()

		case block := <-n.newSyncedBlocks:
			if n.isMining {
				blockHash, _ := block.Hash()
				fmt.Printf("\nPeer mined next Block '%x' faster :(\n", blockHash)

				n.removeMinedPendingTransactions(&block)
				stopCurrentMining()
			}

		case <-ctx.Done():
			ticker.Stop()
			return nil
		}
	}
}

func (n *Node) minePendingTransactions(ctx context.Context) error {
	blockToMine := PendingBlock{
		Previous:     n.state.LatestBlockHash,
		Number:       n.state.LatestBlock.Header.Number + 1,
		Time:         uint64(time.Now().Unix()),
		Transactions: n.getPendingTransactionsArray(),
		Miner:        n.Miner,
	}

	mined, err := Mine(ctx, blockToMine)
	if err != nil {
		return err
	}

	n.removeMinedPendingTransactions(&mined)

	_, err = n.state.AddBlock(mined)
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) removeMinedPendingTransactions(block *munichain.Block) {
	for _, tx := range block.Transactions {
		hash, err := tx.Hash()
		if err != nil {
			continue
		}
		n.archivedTransactions[hash] = tx
		delete(n.pendingTransactions, hash)
	}
}

func (n *Node) addPendingTransaction(tx munichain.Transaction) error {
	txHash, err := tx.Hash()
	if err != nil {
		return err
	}

	_, isAlreadyPending := n.pendingTransactions[txHash]
	_, isArchived := n.archivedTransactions[txHash]

	if !isAlreadyPending || isArchived {
		n.pendingTransactions[txHash] = tx
	}

	return nil
}

func (n *Node) getPendingTransactionsArray() []munichain.Transaction {
	txs := make([]munichain.Transaction, len(n.pendingTransactions))

	i := 0
	for _, tx := range n.pendingTransactions {
		txs[i] = tx
		i++
	}

	return txs
}
