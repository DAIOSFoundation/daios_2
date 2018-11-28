package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
)

type Transactions []*Transaction

type Transaction struct {
	Nonce int  `json:"nonce"`
	Data  Data `json:"data"`
}

type Data struct {
	Hash      string  `json:"hash"`
	Sender    Address `json:"sender"`
	Recipient Address `json:"recipient"`
	Payload   string  `json:"payload"`
	Value     int     `json:"value"`
}

func (txs *Transactions) Len() int { return len(*txs) }

func (txs *Transactions) Push(x *Transaction) {
	*txs = append(*txs, x)
}

func (txs *Transactions) Pop() interface{} {
	old := *txs
	n := len(old)
	x := old[n-1]
	*txs = old[0 : n-1]
	return x
}

/*
func (txs *Transactions) HashTransactions() (*mt.MerkleTree, error) {
	var transactions [][]byte

	for _, tx := range *txs {
		transactions = append(transactions, tx.EncodeTx())
	}
	mTree, err := mt.NewMerkleTree(transactions)
	if err != nil {
		return nil, err
	}

	return mTree, err
}
*/
func (tx *Transaction) Encode() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	encoder.Encode(tx)

	return result.Bytes()
}

func DecodeTx(d []byte) *Transaction {
	var tx *Transaction

	decoder := gob.NewDecoder(bytes.NewReader(d))
	decoder.Decode(&tx)

	return tx
}

func NewTransaction(sender, recipient Address, value int, payload string) *Transaction {
	d := Data{
		Sender:    sender,
		Recipient: recipient,
		Payload:   payload,
		Value:     value,
	}

	return &Transaction{Data: d}
}

func (tx *Transaction) Hash() string {
	h := sha256.New()
	h.Write(tx.Encode())
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func (tx *Transaction) MarshalJSON() []byte {
	hash := tx.Hash()
	data := tx.Data
	data.Hash = hash

	jtx, err := json.Marshal(tx)
	if err != nil {
		panic(err)
	}

	return jtx
}
