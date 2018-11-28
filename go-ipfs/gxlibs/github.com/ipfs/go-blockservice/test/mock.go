package bstest

import (
	. "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-blockservice"

	bitswap "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-bitswap"
	tn "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-bitswap/testnet"
	delay "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-ipfs-delay"
	mockrouting "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-ipfs-routing/mock"
)

// Mocks returns |n| connected mock Blockservices
func Mocks(n int) []BlockService {
	net := tn.VirtualNetwork(mockrouting.NewServer(), delay.Fixed(0))
	sg := bitswap.NewTestSessionGenerator(net)

	instances := sg.Instances(n)

	var servs []BlockService
	for _, i := range instances {
		servs = append(servs, New(i.Blockstore(), i.Exchange))
	}
	return servs
}
