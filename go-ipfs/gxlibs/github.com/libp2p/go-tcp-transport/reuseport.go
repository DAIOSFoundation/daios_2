package tcp

import (
	"os"
	"strings"

	reuseport "github.com/dai/go-ipfs/gxlibs/github.com/libp2p/go-reuseport"
)

// envReuseport is the env variable name used to turn off reuse port.
// It default to true.
const envReuseport = "LIBP2P_TCP_REUSEPORT"
const deprecatedEnvReuseport = "IPFS_REUSEPORT"

// envReuseportVal stores the value of envReuseport. defaults to true.
var envReuseportVal = true

func init() {
	v := strings.ToLower(os.Getenv(envReuseport))
	if v == "false" || v == "f" || v == "0" {
		envReuseportVal = false
		log.Infof("REUSEPORT disabled (LIBP2P_TCP_REUSEPORT=%s)", v)
	}
	v, exist := os.LookupEnv(deprecatedEnvReuseport)
	if exist {
		log.Warning("IPFS_REUSEPORT is deprecated, use LIBP2P_TCP_REUSEPORT instead")
		if v == "false" || v == "f" || v == "0" {
			envReuseportVal = false
			log.Infof("REUSEPORT disabled (IPFS_REUSEPORT=%s)", v)
		}
	}
}

// reuseportIsAvailable returns whether reuseport is available to be used. This
// is here because we want to be able to turn reuseport on and off selectively.
// For now we use an ENV variable, as this handles our pressing need:
//
//   LIBP2P_TCP_REUSEPORT=false ipfs daemon
//
// If this becomes a sought after feature, we could add this to the config.
// In the end, reuseport is a stop-gap.
func ReuseportIsAvailable() bool {
	return envReuseportVal && reuseport.Available()
}
