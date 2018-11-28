package bitswap

import (
	bsnet "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-bitswap/network"
	peer "github.com/dai/go-ipfs/gxlibs/github.com/libp2p/go-libp2p-peer"
	"github.com/dai/go-ipfs/gxlibs/github.com/libp2p/go-testutil"
)

type Network interface {
	Adapter(testutil.Identity) bsnet.BitSwapNetwork

	HasPeer(peer.ID) bool
}
