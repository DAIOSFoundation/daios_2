package ipns

import (
	"context"

	"github.com/dai/go-ipfs/core"
	nsys "github.com/dai/go-ipfs/namesys"
	ci "github.com/dai/go-ipfs/gxlibs/github.com/libp2p/go-libp2p-crypto"
	path "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-path"
	ft "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-unixfs"
)

// InitializeKeyspace sets the ipns record for the given key to
// point to an empty directory.
func InitializeKeyspace(n *core.IpfsNode, key ci.PrivKey) error {
	ctx, cancel := context.WithCancel(n.Context())
	defer cancel()

	emptyDir := ft.EmptyDirNode()

	err := n.Pinning.Pin(ctx, emptyDir, false)
	if err != nil {
		return err
	}

	err = n.Pinning.Flush()
	if err != nil {
		return err
	}

	pub := nsys.NewIpnsPublisher(n.Routing, n.Repo.Datastore())

	return pub.Publish(ctx, key, path.FromCid(emptyDir.Cid()))
}
