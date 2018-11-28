package leveldb

import (
	"testing"

	"github.com/dai/go-ipfs/gxlibs/github.com/syndtr/goleveldb/leveldb/testutil"
)

func TestLevelDB(t *testing.T) {
	testutil.RunSuite(t, "LevelDB Suite")
}
