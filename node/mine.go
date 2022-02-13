package node

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/martinmunillas/munichain/munichain"
)

type PendingBlock struct {
	Previous     munichain.Hash
	Number       uint64
	Time         uint64
	Transactions []munichain.Transaction
	Miner        munichain.Address
}

func Mine(ctx context.Context, pb PendingBlock) (munichain.Block, error) {
	if len(pb.Transactions) == 0 {
		err := fmt.Errorf("mining empty blocks is not allowed")
		return munichain.Block{}, err
	}
	start := time.Now()
	attempt := 0
	var block munichain.Block
	var hash munichain.Hash
	var nonce uint32
	for !munichain.IsBlockHashValid(hash) {
		select {
		case <-ctx.Done():
			fmt.Println("Mining cancelled!")
			err := fmt.Errorf("mining cancelled. %s", ctx.Err())
			return munichain.Block{}, err
		default:
		}
		attempt++
		nonce = generateNonce()
		if attempt%1000000 == 0 || attempt == 1 {
			fmt.Printf("Mining %d Pending TXs. Attempt: %d\n", len(pb.Transactions), attempt)
		}
		block = munichain.Block{
			Header: munichain.BlockHeader{
				Previous: pb.Previous,
				Number:   pb.Number,
				Time:     pb.Time,
				Nonce:    nonce,
				Miner:    pb.Miner,
			},
			Transactions: pb.Transactions,
		}
		blockHash, err := block.Hash()
		if err != nil {
			err = fmt.Errorf("couldn't mine block. %s", err.Error())
			return munichain.Block{}, err
		}
		hash = blockHash
	}
	fmt.Printf("\nMined new Block '%x' using PoW:\n", hash)
	fmt.Printf("\tHeight: '%v'\n", block.Header.Number)
	fmt.Printf("\tNonce: '%v'\n", block.Header.Nonce)
	fmt.Printf("\tCreated: '%v'\n", block.Header.Time)
	fmt.Printf("\tMiner: '%v'\n", block.Header.Miner)
	fmt.Printf("\tPrevious: '%s'\n\n", block.Header.Previous.ToString())
	fmt.Printf("\tAttempt: '%v'\n", attempt)
	fmt.Printf("\tTime: %s\n\n", time.Since(start))
	return block, nil
}

func generateNonce() uint32 {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Uint32()
}
