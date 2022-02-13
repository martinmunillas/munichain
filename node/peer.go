package node

import (
	"fmt"
	"net/http"

	"github.com/martinmunillas/munichain/munichain"
)

type PeerNode struct {
	IP          string `json:"ip"`
	Port        uint64 `json:"port"`
	IsBootstrap bool   `json:"isBootstrap"`
	IsActive    bool   `json:"isActive"`
	connected   bool
}

func (p *PeerNode) TcpAddress() string {
	return fmt.Sprintf("%s:%d", p.IP, p.Port)
}

func (p *PeerNode) fetchBlocksFrom(hash munichain.Hash) ([]munichain.Block, error) {

	res, err := http.Get(
		fmt.Sprintf(
			"http://%s/%s?%s=%s",
			p.TcpAddress(),
			syncEndpoint,
			syncFromBlockQueryKey,
			hash.ToString(),
		),
	)
	if err != nil {
		return nil, err
	}

	blocks := []munichain.Block{}
	err = readRes(res, &blocks)
	if err != nil {
		return nil, err
	}

	return blocks, nil

}

func (p *PeerNode) queryStatus() (StatusRes, error) {
	res, err := http.Get(fmt.Sprintf("http://%s/%s", p.TcpAddress(), statusEndpoint))
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

func (n *Node) removePeer(peer PeerNode) {
	delete(n.KnownPeers, peer.TcpAddress())
}

func (n *Node) addPeer(peer PeerNode) {
	n.KnownPeers[peer.TcpAddress()] = peer
}

func (n *Node) isKnownPeer(peer PeerNode) bool {
	_, ok := n.KnownPeers[peer.TcpAddress()]
	return ok
}
