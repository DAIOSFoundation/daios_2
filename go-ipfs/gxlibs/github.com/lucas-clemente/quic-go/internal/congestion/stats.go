package congestion

import "github.com/dai/go-ipfs/gxlibs/github.com/lucas-clemente/quic-go/internal/protocol"

type connectionStats struct {
	slowstartPacketsLost protocol.PacketNumber
	slowstartBytesLost   protocol.ByteCount
}
