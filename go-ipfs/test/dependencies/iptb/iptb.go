package main

import (
	"fmt"
	"os"

	cli "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/iptb/cli"
	testbed "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/iptb/testbed"

	plugin "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/iptb-plugins/local"
)

func init() {
	_, err := testbed.RegisterPlugin(testbed.IptbPlugin{
		From:        "<builtin>",
		NewNode:     plugin.NewNode,
		GetAttrList: plugin.GetAttrList,
		GetAttrDesc: plugin.GetAttrDesc,
		PluginName:  plugin.PluginName,
		BuiltIn:     true,
	}, false)

	if err != nil {
		panic(err)
	}
}

func main() {
	cli := cli.NewCli()
	if err := cli.Run(os.Args); err != nil {
		fmt.Fprintf(cli.ErrWriter, "%s\n", err)
		os.Exit(1)
	}
}