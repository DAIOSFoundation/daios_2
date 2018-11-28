package iterator_test

import (
	"testing"

	"github.com/dai/go-ipfs/gxlibs/github.com/syndtr/goleveldb/leveldb/testutil"
)

func TestIterator(t *testing.T) {
	testutil.RunSuite(t, "Iterator Suite")
}
