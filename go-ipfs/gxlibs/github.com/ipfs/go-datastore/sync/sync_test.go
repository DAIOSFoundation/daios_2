package sync

import (
	"testing"

	ds "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-datastore"
	dstest "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-datastore/test"
)

func TestSync(t *testing.T) {
	dstest.SubtestAll(t, MutexWrap(ds.NewMapDatastore()))
}
