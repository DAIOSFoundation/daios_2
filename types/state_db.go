package types

import (
	"fmt"
	"math/big"

	"github.com/dai/db"
	"github.com/dai/mt"
)

type StateDB struct {
	db db.BoltDatabase

	tree   *mt.MerkleTree
	States map[interface{}]*State
}

func NewStateDB(root []byte, db db.BoltDatabase) *StateDB {

	var mtree *mt.MerkleTree
	mtree = mt.Decode(db.GetData(root))

	if mtree == nil {
		var data []mt.Data
		data = append(data, mt.Data{Key: []byte(""), Data: []byte("")})

		mtree = mt.NewMerkleTree(data, db)

	}

	return &StateDB{
		db:     db,
		tree:   mtree,
		States: make(map[interface{}]*State),
	}

}

func (sdb *StateDB) CreateAccount(addr Address) *State {
	new, prev := sdb.createState(addr)

	if prev != nil {
		new.SetBalance(prev.Data.Balance)
	}

	return new
}

func (sdb *StateDB) createState(addr Address) (ns, ps *State) {
	ps = sdb.getState(addr)
	ns = NewState(sdb, addr, Account{})
	ns.SetNonce(0)

	sdb.setState(ns)
	return ns, ps
}

func (sdb *StateDB) getState(addr Address) (s *State) {

	if obj := sdb.States[addr]; obj != nil {
		return obj
	}

	enc := sdb.db.GetData(addr[:])

	if len(enc) == 0 {
		return nil
	}

	data := DecodeAccount(enc)

	obj := NewState(sdb, addr, data)
	sdb.setState(obj)
	return obj
}

func (sdb *StateDB) setState(s *State) {
	sdb.States[s.Address] = s
}

func (s *StateDB) GetState(addr Address, bhash []byte) []byte {
	state := s.getStateObject(addr)
	if state != nil {
		return state.GetState(s.db, bhash)
	}
	return []byte{}
}

func (s *StateDB) GetBalance(addr Address) *big.Int {
	state := s.getStateObject(addr)
	if state != nil {
		return state.Balance()
	}
	return big.NewInt(0)
}

func (s *StateDB) AddBalance(addr Address, amount *big.Int) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.AddBalance(amount)

	}
}

func (s *StateDB) SubBalance(addr Address, amount *big.Int) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SubBalance(amount)

	}
}

func (s *StateDB) SetBalance(addr Address, amount *big.Int) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetBalance(amount)

	}
}

func (s *StateDB) SetNonce(addr Address, nonce uint64) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetNonce(nonce)
	}
}

func (s *StateDB) SetState(addr Address, key, value []byte) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetState(s.db, key, value)

	}

}

func (s *StateDB) GetOrNewStateObject(addr Address) *State {
	stateObject := s.getStateObject(addr)
	if stateObject == nil {
		stateObject, _ = s.createState(addr)
	}
	return stateObject
}

func (s *StateDB) getStateObject(addr Address) (stateObject *State) {

	if obj := s.States[addr]; obj != nil {
		return obj
	}

	/*
		enc := s.tree.GetData(addr[:])
		if len(enc) == 0 {
			return nil
		}
	*/
	enc := s.db.GetData(addr[:])

	data := DecodeAccount(enc)

	obj := NewState(s, addr, data)
	s.setStateObject(obj)
	return obj
}

func (s *StateDB) setStateObject(state *State) {
	s.States[state.Address] = state
}

func (s *StateDB) Sync(state map[string]*State) {
	for key, value := range state {
		strKey := fmt.Sprintf("%s", key)
		strValue := value

		s.States[strKey] = strValue

	}
}
