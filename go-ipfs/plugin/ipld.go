package plugin

import (
	"github.com/dai/go-ipfs/core/coredag"

	ipld "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-ipld-format"
)

// PluginIPLD is an interface that can be implemented to add handlers for
// for different IPLD formats
type PluginIPLD interface {
	Plugin

	RegisterBlockDecoders(dec ipld.BlockDecoder) error
	RegisterInputEncParsers(iec coredag.InputEncParsers) error
}
