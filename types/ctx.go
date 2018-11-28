package types

import (
	"context"

	"github.com/dai/go-ipfs/core"
	config "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-ipfs-config"
)

type Context struct {
	Path   string
	Config *config.Config
	Node   *core.IpfsNode
	Ctx    context.Context
}
