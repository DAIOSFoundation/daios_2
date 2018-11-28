package core

import (
	"sync"

	"github.com/dai/db"
	"github.com/dai/types"
)

//const dbBlock = "bc_%s"
const dbBlock = "bc.db"
const blocksBucket = "blocks"

//const dbState = "st_%s"
const dbState = "st.db"
const statesBucket = "states"

type Daios struct {
	blockchain *BlockChain
	txPool     *TxPool
	mutex      *sync.Mutex
	validators []types.Address
}

func (s *Daios) BlockChain() *BlockChain      { return s.blockchain }
func (s *Daios) TxPool() *TxPool              { return s.txPool }
func (s *Daios) Mutex() *sync.Mutex           { return s.mutex }
func (s *Daios) Validators() *[]types.Address { return &s.validators }

func New(nodeID types.Address) *Daios {
	Daios := Daios{}

	//dbBlock := fmt.Sprintf(dbBlock, nodeID)
	blockDB := db.NewDataBase(dbBlock, blocksBucket)

	defer blockDB.Close()

	data := blockDB.GetData([]byte("l"))

	if data == nil {
		Daios.blockchain = NewBlockChain(*blockDB, nil)
		genesisBlock := Daios.blockchain.GetgenesisBlock()
		types.NewBlock(genesisBlock, nil, types.Address{}, nil)

		blockDB.SetData([]byte(genesisBlock.Hash), genesisBlock.Encode())
		blockDB.SetData([]byte("l"), []byte(genesisBlock.Hash))

	} else {
		blockHash := blockDB.GetData([]byte("l"))
		block := blockDB.GetData([]byte(blockHash))
		genesisblock := types.Decode(block)

		Daios.blockchain = NewBlockChain(*blockDB, genesisblock)
	}

	//dbState := fmt.Sprintf(dbState, nodeID)
	stateDB := db.NewDataBase(dbState, statesBucket)
	defer stateDB.Close()

	sdb := types.NewStateDB([]byte("A"), *stateDB)
	types.CS = sdb.CreateAccount(types.NewAddress("ED968E840D"))
	types.CS = sdb.CreateAccount(types.NewAddress("61C5F9D848"))
	types.CS = sdb.CreateAccount(types.NewAddress("EBDD83C1B1"))

	Daios.txPool = NewTxPool(*Daios.blockchain)
	Daios.mutex = &sync.Mutex{}

	return &Daios
}
