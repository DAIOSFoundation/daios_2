package coremock

import (
	"context"

	commands "github.com/dai/go-ipfs/commands"
	core "github.com/dai/go-ipfs/core"
	"github.com/dai/go-ipfs/repo"

	pstore "github.com/dai/go-ipfs/gxlibs/github.com/libp2p/go-libp2p-peerstore"
	host "github.com/dai/go-ipfs/gxlibs/github.com/libp2p/go-libp2p-host"
	libp2p "github.com/dai/go-ipfs/gxlibs/github.com/libp2p/go-libp2p"
	mocknet "github.com/dai/go-ipfs/gxlibs/github.com/libp2p/go-libp2p/p2p/net/mock"
	testutil "github.com/dai/go-ipfs/gxlibs/github.com/libp2p/go-testutil"
	datastore "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-datastore"
	syncds "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-datastore/sync"
	config "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-ipfs-config"
	peer "github.com/dai/go-ipfs/gxlibs/github.com/libp2p/go-libp2p-peer"
)

// NewMockNode constructs an IpfsNode for use in tests.
func NewMockNode() (*core.IpfsNode, error) {
	ctx := context.Background()

	// effectively offline, only peer in its network
	return core.NewNode(ctx, &core.BuildCfg{
		Online: true,
		Host:   MockHostOption(mocknet.New(ctx)),
	})
}

func MockHostOption(mn mocknet.Mocknet) core.HostOption {
	return func(ctx context.Context, id peer.ID, ps pstore.Peerstore, _ ...libp2p.Option) (host.Host, error) {
		return mn.AddPeerWithPeerstore(id, ps)
	}
}

func MockCmdsCtx() (commands.Context, error) {
	// Generate Identity
	ident, err := testutil.RandIdentity()
	if err != nil {
		return commands.Context{}, err
	}
	p := ident.ID()

	conf := config.Config{
		Identity: config.Identity{
			PeerID: p.String(),
		},
	}

	r := &repo.Mock{
		D: syncds.MutexWrap(datastore.NewMapDatastore()),
		C: conf,
	}

	node, err := core.NewNode(context.Background(), &core.BuildCfg{
		Repo: r,
	})
	if err != nil {
		return commands.Context{}, err
	}

	return commands.Context{
		Online:     true,
		ConfigRoot: "/tmp/.mockipfsconfig",
		LoadConfig: func(path string) (*config.Config, error) {
			return &conf, nil
		},
		ConstructNode: func() (*core.IpfsNode, error) {
			return node, nil
		},
	}, nil
}
