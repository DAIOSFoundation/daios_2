package pnet

import (
	"fmt"
	"io"

	mc "github.com/dai/go-ipfs/gxlibs/github.com/multiformats/go-multicodec"
	bmux "github.com/dai/go-ipfs/gxlibs/github.com/multiformats/go-multicodec/base/mux"
)

var (
	pathPSKv1   = []byte("/key/swarm/psk/1.0.0/")
	headerPSKv1 = mc.Header(pathPSKv1)
)

func decodeV1PSK(in io.Reader) (*[32]byte, error) {
	var err error
	in, err = mc.WrapTransformPathToHeader(in)
	if err != nil {
		return nil, err
	}
	err = mc.ConsumeHeader(in, headerPSKv1)
	if err != nil {
		return nil, fmt.Errorf("psk header error: %s", err.Error())
	}

	in, err = mc.WrapTransformPathToHeader(in)
	if err != nil {
		return nil, fmt.Errorf("wrapping error: %s", err.Error())
	}
	out := [32]byte{}

	err = bmux.AllBasesMux().Decoder(in).Decode(out[:])
	return &out, err
}
