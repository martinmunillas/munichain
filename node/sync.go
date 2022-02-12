package node

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func (n *Node) sync(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 45)

	for {
		select {
		case <-ticker.C:
			fmt.Println("Syncing...")

			n.fetchNewBlocksAndPeers()
		case <-ctx.Done():
			ticker.Stop()
		}
	}
}

func (n *Node) fetchNewBlocksAndPeers() {
	for _, peer := range n.KnownPeers {
		status, err := queryPeerStatus(peer)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		localBlockNumber := n.state.LatestBlock.Header.Number
		if localBlockNumber < status.Number {
			newBlocksCount := status.Number - localBlockNumber
			fmt.Printf("Found %d new blocks from Peer %s\n", newBlocksCount, peer.IP)
		}

		for _, statusPeer := range status.KnownPeers {
			newPeer, isKnownPeer := n.KnownPeers[statusPeer.TcpAddress()]
			if !isKnownPeer {
				fmt.Printf("New peer: %s\n", statusPeer.TcpAddress())

				n.KnownPeers[statusPeer.TcpAddress()] = newPeer
			}
		}
	}
}

func queryPeerStatus(peer PeerNode) (StatusRes, error) {
	url := fmt.Sprintf("http://%s/%s", peer.TcpAddress(), statusEndpoint)
	res, err := http.Get(url)
	if err != nil {
		return StatusRes{}, err
	}

	statusRes := StatusRes{}
	err = readRes(res, &statusRes)
	if err != nil {
		return StatusRes{}, err
	}

	return statusRes, nil
}
