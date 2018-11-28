package db

import (
	"log"
	"os"
	"strconv"

	"github.com/boltdb/bolt"
)

type BoltDatabase struct {
	database     *bolt.DB
	fileName     string
	blocksBucket string
}

func NewDataBase(fn string, bb string) *BoltDatabase {
	/*
		if exists(fn) {
			fmt.Println("db file exists.")
			os.Exit(1)
		}
	*/

	bdb, err := bolt.Open(fn, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &BoltDatabase{fileName: fn, database: bdb, blocksBucket: bb}
}

func (db *BoltDatabase) Close() {
	db.database.Close()
}

func (db *BoltDatabase) FileName() string {
	return db.fileName
}
func (db *BoltDatabase) BoltDB() *bolt.DB {
	return db.database
}

func exists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func (db *BoltDatabase) SetData(key []byte, value []byte) error {
	return db.database.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(db.blocksBucket))
		if err != nil {
			return err
		}

		return bkt.Put(key, value)
	})
}

func (db *BoltDatabase) GetLastKey() (uint, error) {
	var id uint64
	err := db.database.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(db.blocksBucket))
		if err != nil {
			panic(err)
		}
		err = nil

		c := bkt.Cursor()
		k, _ := c.Last()
		id, err = strconv.ParseUint(string(k), 16, 64)

		return err
	})

	return uint(id), err
}

func (db *BoltDatabase) GetData(id []byte) []byte {
	var data []byte
	db.database.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(db.blocksBucket))
		if bkt == nil {
			return nil
		}

		data = bkt.Get(id)

		return nil

	})
	return data
}

func (db *BoltDatabase) GetBuckets() ([]string, error) {
	var names = make([]string, 0)
	err := db.database.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			names = append(names, string(name))
			return nil
		})
		return nil
	})

	return names, err
}
