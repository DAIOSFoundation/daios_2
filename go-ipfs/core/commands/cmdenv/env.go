package cmdenv

import (
	"fmt"

	"github.com/dai/go-ipfs/commands"
	"github.com/dai/go-ipfs/core"
	coreiface "github.com/dai/go-ipfs/core/coreapi/interface"

	cmds "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-ipfs-cmds"
	config "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-ipfs-config"
)

// GetNode extracts the node from the environment.
func GetNode(env interface{}) (*core.IpfsNode, error) {
	ctx, ok := env.(*commands.Context)
	if !ok {
		return nil, fmt.Errorf("expected env to be of type %T, got %T", ctx, env)
	}

	return ctx.GetNode()
}

// GetApi extracts CoreAPI instance from the environment.
func GetApi(env cmds.Environment) (coreiface.CoreAPI, error) {
	ctx, ok := env.(*commands.Context)
	if !ok {
		return nil, fmt.Errorf("expected env to be of type %T, got %T", ctx, env)
	}

	return ctx.GetApi()
}

// GetConfig extracts the config from the environment.
func GetConfig(env cmds.Environment) (*config.Config, error) {
	ctx, ok := env.(*commands.Context)
	if !ok {
		return nil, fmt.Errorf("expected env to be of type %T, got %T", ctx, env)
	}

	return ctx.GetConfig()
}

// GetConfigRoot extracts the config root from the environment
func GetConfigRoot(env cmds.Environment) (string, error) {
	ctx, ok := env.(*commands.Context)
	if !ok {
		return "", fmt.Errorf("expected env to be of type %T, got %T", ctx, env)
	}

	return ctx.ConfigRoot, nil
}
