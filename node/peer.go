package node

import "fmt"

type PeerNode struct {
	IP          string `json:"ip"`
	Port        uint64 `json:"port"`
	IsBootstrap bool   `json:"isBootstrap"`
	IsActive    bool   `json:"isActive"`
}

func (p *PeerNode) TcpAddress() string {
	return fmt.Sprintf("%s:%d", p.IP, p.Port)
}
