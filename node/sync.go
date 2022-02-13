package node

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/martinmunillas/munichain/munichain"
)

func (n *Node) sync(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 10)

	n.doSync()

	for {
		select {
		case <-ticker.C:
			n.doSync()
		case <-ctx.Done():
			ticker.Stop()
		}
	}
}

func (n *Node) doSync() {
	fmt.Println("Syncing...")
	for _, peer := range n.KnownPeers {
		status, err := peer.queryStatus()
		if err != nil {

			fmt.Printf("Error: %s\n", err)
			fmt.Printf("Peer '%s' is being removed from the network", peer.TcpAddress())

			n.removePeer(peer)
			continue
		}

		err = n.joinKnownPeers(peer)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		err = n.syncBlocks(peer, status)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		err = n.syncKnownPeers(status)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		err = n.syncPendingTransactions(status.PendingTransactions)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
	}
}

func (n *Node) joinKnownPeers(peer PeerNode) error {
	if peer.connected {
		return nil
	}

	url := fmt.Sprintf(
		"http://%s%s?%s=%s&%s=%d",
		peer.TcpAddress(),
		joinPeerEndpoint,
		joinPeerIPQueryKey,
		n.IP,
		joinPeerPortQueryKey,
		n.Port,
	)

	res, err := http.Get(url)
	if err != nil {
		return err
	}

	addPeerRes := AddPeerRes{}
	err = readRes(res, &addPeerRes)
	if err != nil {
		return err
	}

	knownPeer := n.KnownPeers[peer.TcpAddress()]
	knownPeer.connected = addPeerRes.Success

	n.addPeer(knownPeer)

	if !addPeerRes.Success {
		return fmt.Errorf("unable to join KnownPeers of '%s'", peer.TcpAddress())
	}

	return nil

}

func (n *Node) syncBlocks(peer PeerNode, status StatusRes) error {
	localBlockNumber := n.state.LatestBlock.Header.Number
	if status.Number <= localBlockNumber {
		return nil
	}

	newBlocksAmount := status.Number - localBlockNumber
	var blockStr string
	if newBlocksAmount == 1 {
		blockStr = "block"
	} else {
		blockStr = "blocks"
	}
	fmt.Printf("Found %d new %s from Peer %s\n", newBlocksAmount, blockStr, peer.TcpAddress())

	blocks, err := peer.fetchBlocksFrom(n.state.LatestBlockHash)
	if err != nil {
		return err
	}
	for _, block := range blocks {
		_, err = n.state.AddBlock(block)
		if err != nil {
			return err
		}

		n.newSyncedBlocks <- block
	}
	return nil
}

func (n *Node) syncKnownPeers(status StatusRes) error {
	for _, statusPeer := range status.KnownPeers {
		if !n.isKnownPeer(statusPeer) {
			fmt.Printf("found new peer: %s\n", statusPeer.TcpAddress())

			n.addPeer(statusPeer)
		}
	}
	return nil
}

func (n *Node) syncPendingTransactions(txs []munichain.Transaction) error {
	for _, tx := range txs {
		err := n.addPendingTransaction(tx)
		if err != nil {
			return err
		}
	}

	return nil
}
