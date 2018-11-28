package memdb

import (
	"testing"

	"github.com/dai/go-ipfs/gxlibs/github.com/syndtr/goleveldb/leveldb/testutil"
)

func TestMemDB(t *testing.T) {
	testutil.RunSuite(t, "MemDB Suite")
}
