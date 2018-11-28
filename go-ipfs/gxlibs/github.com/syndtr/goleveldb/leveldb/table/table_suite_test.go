package table

import (
	"testing"

	"github.com/dai/go-ipfs/gxlibs/github.com/syndtr/goleveldb/leveldb/testutil"
)

func TestTable(t *testing.T) {
	testutil.RunSuite(t, "Table Suite")
}
