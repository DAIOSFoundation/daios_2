package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"time"

	"github.com/dai/core"
	datastore "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-datastore"
	net "github.com/dai/go-ipfs/gxlibs/github.com/libp2p/go-libp2p-net"

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

	host, err := libp2p.New(ctx)
	if err != nil {
		panic(err)
	}

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

	host.SetStreamHandler("/ipfs/id/1.0.0", handleStream)

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
func handleStream(s net.Stream) {

	fmt.Println("Connected :" + s.Conn().RemotePeer().Pretty())
	list := make(map[string]types.Address)
	addr := types.NewAddress(s.Conn().RemotePeer().Pretty())
	list[string(addr[:])] = addr
	core.ValidatorPool = append(core.ValidatorPool, addr)

}
