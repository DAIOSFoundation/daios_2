package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"time"

	"github.com/dai/core"
	datastore "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-datastore"
	pnet "github.com/dai/go-ipfs/gxlibs/github.com/libp2p/go-libp2p-pnet"

	//"github.com/libp2p/go-libp2p/config"

	//"github.com/dai/go-ipfs/gxlibs/github.com/libp2p/go-floodsub"
	"github.com/dai/go-ipfs/gxlibs/github.com/libp2p/go-libp2p"
	libp2pdht "github.com/dai/go-ipfs/gxlibs/github.com/libp2p/go-libp2p-kad-dht"
	floodsub "github.com/dai/go-ipfs/gxlibs/github.com/libp2p/go-libp2p-pubsub"
	"github.com/dai/types"
)

var wg sync.WaitGroup

func main() {
	wg.Add(1)
	mineTopic := "daiosMine"
	peerTopic := "daiosPeer"
	list := make(map[string]types.Address)

	ctx := context.Background()
	core.ValidatorPool = []types.Address{}

	opts := []libp2p.Option{}
	swarmkey, err := SwarmKey()
	if swarmkey != nil {
		protec, err := pnet.NewProtector(bytes.NewReader(swarmkey))
		if err != nil {
			fmt.Errorf("failed NewProtector %s", err)
		}
		opts = append(opts, libp2p.PrivateNetwork(protec))
	}

	host, err := libp2p.New(ctx, opts...)
	if err != nil {
		panic(err)
	}
	libp2p.ChainOptions(opts...)
	for _, addr := range host.Addrs() {
		fmt.Printf("%s/ipfs/%s\n", addr.String(), host.ID().Pretty())
	}

	fsub, err := floodsub.NewFloodSub(ctx, host)
	if err != nil {
		panic(err)
	}

	dht := libp2pdht.NewDHT(ctx, host, datastore.NewMapDatastore())
	if err != nil {
		panic(err)
	}

	bsConfig := libp2pdht.DefaultBootstrapConfig
	bsConfig.Period = 10 * time.Second
	bsConfig.Queries = 1000
	if _, err := dht.BootstrapWithConfig(bsConfig); err != nil {
		panic(err)
	}

	go func() {
		for {
			core.Pick()

		}
	}()

	go func() {
		for {
			broadCastAddr := core.BroadCastAddr()
			b := <-*broadCastAddr
			fmt.Println("PICK", b)
			jsonBytes, err := json.Marshal(b)
			if err != nil {
				panic(err)
			}
			if err := fsub.Publish(mineTopic, jsonBytes); err != nil {
				panic(err)
			}
		}
	}()

	subPeer, err := fsub.Subscribe(peerTopic)
	if err != nil {
		panic(err)
	}

	go func() {
		for {

			msg, err := subPeer.Next(ctx)
			if err != nil {
				panic(err)
			}

			var addr types.Address
			json.Unmarshal(msg.GetData(), &addr)

			list[string(addr[:])] = addr

			core.ValidatorPool = append(core.ValidatorPool, addr)

			time.Sleep(1 * time.Second)

		}

	}()
	wg.Wait()
}

func SwarmKey() ([]byte, error) {
	curr, _ := os.Getwd()
	spath := filepath.Join(curr, "swarm.key")
	fmt.Println(spath)
	f, err := os.Open(spath)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return nil, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}
