package types

import (
	"bytes"
	"encoding/gob"
	"math/big"

	"github.com/dai/db"
	"github.com/dai/mt"
)

type Account struct {
	Nonce          uint64
	Balance        *big.Int
	MerkleRootHash string
}

type State struct {
	AddrHash Hash
	Address  Address
	Data     Account

	db   *StateDB
	tree mt.MerkleTree
}

var CS *State

func (s *State) empty() bool {
	return s.Data.Nonce == 0 && s.Data.Balance.Sign() == 0
}

func (s *State) Balance() *big.Int {
	return s.Data.Balance
}

func (s *State) DB() *StateDB {
	return s.db
}

func NewState(db *StateDB, address Address, data Account) *State {
	if data.Balance == nil {
		data.Balance = new(big.Int)
	}

	return &State{
		db:       db,
		Address:  address,
		AddrHash: Hash(address[:]),
		Data:     data,
	}
}

func (s *State) AddBalance(amount *big.Int) {
	if amount.Sign() == 0 {
		return
	}
	s.SetBalance(new(big.Int).Add(s.Balance(), amount))
}

func (s *State) SubBalance(amount *big.Int) {
	if amount.Sign() == 0 {
		return
	}
	s.SetBalance(new(big.Int).Sub(s.Balance(), amount))
}

func (s *State) SetBalance(amount *big.Int) {
	s.Data.Balance = amount
}

func (s *State) SetNonce(nonce uint64) {
	s.Data.Nonce = nonce
}

func (s *State) GetState(db db.BoltDatabase, key []byte) []byte {

	// object를 db로 부터 읽는다
	//enc := s.tree.GetData([]byte(key))

	enc := db.GetData(key)
	if len(enc) == 0 {
		return nil
	}

	return enc
}

func (s *State) SetState(db db.BoltDatabase, key, value []byte) {

	prev := s.GetState(db, key)

	if bytes.Compare(prev, []byte(value[:])) == 0 {
		return
	}

	db.SetData(key, value)

}

func (ac *Account) Encode() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	encoder.Encode(ac)

	return result.Bytes()
}

func (s *State) Encode() []byte {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	encoder.Encode(s.Data)

	return result.Bytes()

}

func DecodeAccount(d []byte) Account {
	var a Account

	decoder := gob.NewDecoder(bytes.NewReader(d))
	decoder.Decode(&a)

	return a
}

func StateDecode(d []byte) map[string]*State {
	var ac map[string]*State

	decoder := gob.NewDecoder(bytes.NewReader(d))

	decoder.Decode(ac)

	return ac
}

func StateEncode(s map[interface{}]*State) []byte {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	encoder.Encode(s)

	return result.Bytes()
}
