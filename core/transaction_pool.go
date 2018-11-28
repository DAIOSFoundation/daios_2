package core

import (
	"github.com/dai/types"
)

type TxPool struct {
	chain   BlockChain
	queue   *types.Transactions
	pending *types.Transactions
}

func (txp *TxPool) Enqueue(tx *types.Transaction) {
	*txp.queue = append(*txp.queue, tx)

}

func (txp *TxPool) Dequeue() *types.Transaction { //interface{} {
	old := *txp.queue
	n := len(old)
	x := old[n-1]
	*txp.queue = old[0 : n-1]
	return x
}

func NewTxPool(bc BlockChain) *TxPool {
	txPool := &TxPool{
		chain: bc,
		queue: &types.Transactions{},
	}
	return txPool
}

func (txp *TxPool) Remove() {
	txp.queue = &types.Transactions{}
}
