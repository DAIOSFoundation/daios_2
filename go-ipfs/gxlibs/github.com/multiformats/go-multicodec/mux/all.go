package muxcodec

import (
	mc "github.com/dai/go-ipfs/gxlibs/github.com/multiformats/go-multicodec"
	cbor "github.com/dai/go-ipfs/gxlibs/github.com/multiformats/go-multicodec/cbor"
	json "github.com/dai/go-ipfs/gxlibs/github.com/multiformats/go-multicodec/json"
)

func StandardMux() *Multicodec {
	return MuxMulticodec([]mc.Multicodec{
		cbor.Multicodec(),
		json.Multicodec(false),
		json.Multicodec(true),
	}, SelectFirst)
}
