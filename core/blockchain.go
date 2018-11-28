package core

import (
	"encoding/json"

	"github.com/dai/db"
	"github.com/dai/types"
)

type BlockChain struct {
	genesisBlock *types.Block

	Blocks []*types.Block
	db     db.BoltDatabase
}

type BlockchainIterator struct {
	BlockHash string
	db        db.BoltDatabase
}

func NewBlockChain(bdb db.BoltDatabase, gb *types.Block) *BlockChain {

	bc := BlockChain{}
	bc.db = bdb
	if gb == nil {
		bc.genesisBlock = types.NewGenesisBlock()
	} else {
		bc.genesisBlock = gb
	}

	bc.Blocks = append(bc.Blocks, bc.genesisBlock)
	return &bc
}

func (bc *BlockChain) AddBlock(nb *types.Block) {

	block := bc.db.GetData([]byte("l"))
	if block != nil {
		nb = types.Decode(block)
	}

	if len(bc.Blocks) == 0 {
		bc.Blocks = append(bc.Blocks, nb)
	} else if types.IsValid(nb, bc.Blocks[len(bc.Blocks)-1]) {
		nbc := append(bc.Blocks, nb)
		replaceChain(&bc.Blocks, &nbc)
	}

	bc.db.SetData([]byte(nb.Hash), nb.Encode())
	bc.db.SetData([]byte("l"), []byte(nb.Hash))
}

func (bc *BlockChain) GetLastBlock() *types.Block { return bc.Blocks[len(bc.Blocks)-1] }

func (bc *BlockChain) GetgenesisBlock() *types.Block { return bc.genesisBlock }

func replaceChain(obc, nbc *[]*types.Block) {
	if len(*nbc) > len(*obc) {
		*obc = *nbc
	}
}

func (bc *BlockChain) SyncChain(nbc *BlockChain) {
	if len(nbc.Blocks) >= len(bc.Blocks) {
		bc.genesisBlock = nbc.genesisBlock
		bc.Blocks = nbc.Blocks
	}
}

func (bc *BlockChain) MarshalJSON() []byte {
	data, err := json.Marshal(bc)
	if err != nil {
		panic(err)
	}

	return data
}

func (bc *BlockChain) Iterator() *BlockchainIterator {
	var block *types.Block

	block = bc.GetLastBlock()

	bci := &BlockchainIterator{block.Hash, bc.db}
	return bci
}

func (bi *BlockchainIterator) Next() *types.Block {
	var block *types.Block

	encodedBlock := bi.db.GetData([]byte(bi.BlockHash))
	block = types.Decode(encodedBlock)

	bi.BlockHash = block.Header.PrevHash
	return block
}
