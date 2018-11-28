package core

import (
	"math/big"
	"math/rand"
	"time"

	"github.com/dai/types"
)

var broadCastAddr = make(chan types.Address)

var ValidatorPool []types.Address

var broadCastBlocks = make(chan *types.Block)

func BroadCastBlocks() *chan *types.Block { return &broadCastBlocks }

func BroadCastAddr() *chan types.Address { return &broadCastAddr }

func Pick() {
	time.Sleep(3 * time.Second)
	if len(ValidatorPool) > 0 {
		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)

		Minner := ValidatorPool[r.Intn(len(ValidatorPool))]
		broadCastAddr <- Minner
	}

}

func Mine(daios *Daios, Address types.Address) {
	mutex := *daios.Mutex()

	if len(*daios.txPool.queue) > 0 {

		mutex.Lock()
		sdb := *types.CS.DB()

		var txQue types.Transactions
		txQue = *daios.txPool.queue

		/*
			for len(*daios.txPool.queue) > 0 {
				tx := *daios.txPool.Dequeue()

				sdb.AddBalance(types.NewAddress(tx.Data.Payload.Wallet), big.NewInt(1))
			}
		*/

		sdb.AddBalance(Address, big.NewInt(10))

		b := types.NewBlock(daios.blockchain.GetLastBlock(), txQue, Address, sdb.States) //types.StateEncode(sdb.States))
		daios.blockchain.AddBlock(b)

		mutex.Unlock()

		broadCastBlocks <- b
	}

}
