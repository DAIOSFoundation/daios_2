package types

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dai/com"
)

type Header struct {
	Number    int     `json:"number"`
	Timestamp string  `json:"timestamp"`
	PrevHash  string  `json:"prevHash"`
	Validator Address `json:"validator"`
}

type Block struct {
	Hash         string            `json:"hash"`
	Header       *Header           `json:"header"`
	Transactions *Transactions     `json:"transactions"`
	States       map[string]*State `json:"state"`
}

func NewGenesisBlock() *Block {
	nb := Block{}

	t := time.Now()

	header := &Header{
		Number:    0,
		Timestamp: t.String(),
		PrevHash:  "",
		Validator: Address{},
	}

	nb.Header = header
	nb.Transactions = nil
	nb.States = nil
	nb.Hash = bHash(&nb)

	return &nb
}

func NewBlock(pb *Block, txs Transactions, validator Address, states map[interface{}]*State) *Block {
	nb := Block{}

	t := time.Now()

	header := &Header{
		Number:    pb.Header.Number + 1,
		Timestamp: t.String(),
		PrevHash:  pb.Hash,
		Validator: validator,
	}

	mapString := make(map[string]*State)

	for key, value := range states {
		strKey := fmt.Sprintf("%v", key)
		strValue := value

		mapString[strKey] = strValue
	}

	nb.Header = header
	nb.Transactions = &txs
	nb.States = mapString
	nb.Hash = bHash(&nb)

	return &nb
}

func (b *Block) Encode() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	encoder.Encode(b)

	return result.Bytes()
}

func Decode(d []byte) *Block {
	var block *Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	decoder.Decode(&block)

	return block
}

func IsValid(nb, pb *Block) bool {

	if pb.Header.Number+1 != nb.Header.Number {
		return false
	}
	if pb.Hash != nb.Header.PrevHash {
		return false
	}
	if bHash(nb) != nb.Hash {
		return false
	}
	return true
}

func (b *Block) MarshalJSON() []byte {

	data, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}

	return data
}

func (h *Header) MarshalJSON() []byte {
	data, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}

	return data
}

func bHash(block *Block) string {
	record := string(block.Header.Number) + block.Header.Timestamp + block.Header.PrevHash
	return com.Hash(record)

}
