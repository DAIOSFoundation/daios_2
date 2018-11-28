package mt

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/dai/com"

	"github.com/dai/db"
)

type Data struct {
	Key  []byte
	Data []byte
}
type MerkleTree struct {
	Root       *Node
	merkleRoot []byte
	Leafs      []*Node
	DB         db.BoltDatabase
}

type Node struct {
	Parent *Node
	Left   *Node
	Right  *Node
	leaf   bool
	dup    bool
	Hash   []byte
	Data   Data
}

func (n *Node) checkNode() ([]byte, error) {
	if n.leaf {
		return n.Hash, nil
	}
	rightBytes, err := n.Right.checkNode()
	if err != nil {
		return nil, err
	}

	leftBytes, err := n.Left.checkNode()
	if err != nil {
		return nil, err
	}

	h := sha256.New()
	if _, err := h.Write(append(leftBytes, rightBytes...)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func (n *Node) nodeHash() ([]byte, error) {
	if n.leaf {
		return n.Hash, nil
	}

	h := sha256.New()
	if _, err := h.Write(append(n.Left.Hash, n.Right.Hash...)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func NewMerkleTree(data []Data, db db.BoltDatabase) *MerkleTree {
	root, leafs, err := build(data)

	if err != nil {
		return nil
	}

	t := &MerkleTree{
		Root:       root,
		DB:         db,
		merkleRoot: root.Hash,
		Leafs:      leafs,
	}
	return t
}

func build(data []Data) (*Node, []*Node, error) {

	var leafs []*Node
	for _, d := range data {

		hash := com.HashByte(d.Key)

		leafs = append(leafs, &Node{
			Hash: hash,
			Data: d,
			leaf: true,
		})
	}
	if len(leafs)%2 == 1 {
		duplicate := &Node{
			Hash: leafs[len(leafs)-1].Hash,
			Data: leafs[len(leafs)-1].Data,
			leaf: true,
			dup:  true,
		}
		leafs = append(leafs, duplicate)
	}
	root, err := intermediate(leafs)
	if err != nil {
		return nil, nil, err
	}

	return root, leafs, nil
}

func intermediate(nl []*Node) (*Node, error) {
	var nodes []*Node
	for i := 0; i < len(nl); i += 2 {
		h := sha256.New()
		var left, right int = i, i + 1
		if i+1 == len(nl) {
			right = i
		}
		chash := append(nl[left].Hash, nl[right].Hash...)
		if _, err := h.Write(chash); err != nil {
			return nil, err
		}
		n := &Node{
			Left:  nl[left],
			Right: nl[right],
			Hash:  h.Sum(nil),
		}
		nodes = append(nodes, n)
		nl[left].Parent = n
		nl[right].Parent = n
		if len(nl) == 2 {
			return n, nil
		}
	}
	return intermediate(nodes)
}

func (m *MerkleTree) MerkleRoot() []byte {
	return m.merkleRoot
}

func (m *MerkleTree) SetDB() []byte {
	log.Println(m.MerkleRoot())
	m.DB.SetData(m.MerkleRoot(), nil)
	return m.MerkleRoot()
}

func (m *MerkleTree) RebuildTreeWith(data []Data, key [][]byte) error {
	root, leafs, err := build(data)
	if err != nil {
		return err
	}
	m.Root = root
	m.Leafs = leafs
	m.merkleRoot = root.Hash
	return nil
}

func (m *MerkleTree) CheckKey(key []byte) (bool, error) {
	for _, l := range m.Leafs {
		ok := bytes.Equal(l.Data.Key, key)

		if ok {
			currentParent := l.Parent
			for currentParent != nil {
				h := sha256.New()
				rightBytes, err := currentParent.Right.nodeHash()
				if err != nil {
					return false, err
				}

				leftBytes, err := currentParent.Left.nodeHash()
				if err != nil {
					return false, err
				}
				if currentParent.Left.leaf && currentParent.Right.leaf {
					if _, err := h.Write(append(leftBytes, rightBytes...)); err != nil {
						return false, err
					}
					if bytes.Compare(h.Sum(nil), currentParent.Hash) != 0 {
						return false, nil
					}
					currentParent = currentParent.Parent
				} else {
					if _, err := h.Write(append(leftBytes, rightBytes...)); err != nil {
						return false, err
					}
					if bytes.Compare(h.Sum(nil), currentParent.Hash) != 0 {
						return false, nil
					}
					currentParent = currentParent.Parent
				}
			}
			return true, nil
		}
	}
	return false, nil
}

func (m *MerkleTree) CheckData(Data []byte) (bool, error) {
	for _, l := range m.Leafs {
		ok := bytes.Equal(l.Data.Data, Data)

		if ok {
			currentParent := l.Parent
			for currentParent != nil {
				h := sha256.New()
				rightBytes, err := currentParent.Right.nodeHash()
				if err != nil {
					return false, err
				}

				leftBytes, err := currentParent.Left.nodeHash()
				if err != nil {
					return false, err
				}
				if currentParent.Left.leaf && currentParent.Right.leaf {
					if _, err := h.Write(append(leftBytes, rightBytes...)); err != nil {
						return false, err
					}
					if bytes.Compare(h.Sum(nil), currentParent.Hash) != 0 {
						return false, nil
					}
					currentParent = currentParent.Parent
				} else {
					if _, err := h.Write(append(leftBytes, rightBytes...)); err != nil {
						return false, err
					}
					if bytes.Compare(h.Sum(nil), currentParent.Hash) != 0 {
						return false, nil
					}
					currentParent = currentParent.Parent
				}
			}
			return true, nil
		}
	}
	return false, nil
}

func (m *MerkleTree) GetData(key []byte) []byte {
	for _, l := range m.Leafs {
		bKey := bytes.Equal(l.Data.Key, key)

		if bKey {
			currentParent := l.Parent
			for currentParent != nil {
				h := sha256.New()
				rightBytes, err := currentParent.Right.nodeHash()
				if err != nil {
					return nil
				}

				leftBytes, err := currentParent.Left.nodeHash()
				if err != nil {
					return nil
				}
				if currentParent.Left.leaf && currentParent.Right.leaf {
					if _, err := h.Write(append(leftBytes, rightBytes...)); err != nil {
						return nil
					}
					if bytes.Compare(h.Sum(nil), currentParent.Hash) != 0 {
						return nil
					}
					currentParent = currentParent.Parent
				} else {
					if _, err := h.Write(append(leftBytes, rightBytes...)); err != nil {
						return nil
					}
					if bytes.Compare(h.Sum(nil), currentParent.Hash) != 0 {
						return nil
					}
					currentParent = currentParent.Parent
				}
			}
			return l.Data.Data
		}
	}
	return nil
}

func (m *MerkleTree) String() string {
	s := ""
	for _, l := range m.Leafs {
		s += fmt.Sprint(l)
		s += "\n"
	}
	return s
}

func (m *MerkleTree) Encode() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	encoder.Encode(m)

	return result.Bytes()
}

func Decode(d []byte) *MerkleTree {
	var mt *MerkleTree

	decoder := gob.NewDecoder(bytes.NewReader(d))
	decoder.Decode(&mt)

	return mt
}
