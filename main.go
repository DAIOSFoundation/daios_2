package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/dai/com"
	dcore "github.com/dai/core"
	version "github.com/dai/go-ipfs"
	"github.com/dai/go-ipfs/core"
	ls "github.com/dai/go-ipfs/core/commands"
	"github.com/dai/go-ipfs/core/coreapi"
	iface "github.com/dai/go-ipfs/core/coreapi/interface"
	corehttp "github.com/dai/go-ipfs/core/corehttp"
	"github.com/dai/go-ipfs/core/corerepo"
	"github.com/dai/go-ipfs/core/coreunix"
	humanize "github.com/dai/go-ipfs/gxlibs/github.com/dustin/go-humanize"
	cid "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-cid"
	cmdsHttp "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-ipfs-cmds/http"
	config "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-ipfs-config"
	ipld "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-ipld-format"
	merkledag "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-merkledag"
	unixfs "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-unixfs"
	uio "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-unixfs/io"
	unixfspb "github.com/dai/go-ipfs/gxlibs/github.com/ipfs/go-unixfs/pb"
	metrics "github.com/dai/go-ipfs/gxlibs/github.com/libp2p/go-libp2p-metrics"
	inet "github.com/dai/go-ipfs/gxlibs/github.com/libp2p/go-libp2p-net"
	ma "github.com/dai/go-ipfs/gxlibs/github.com/multiformats/go-multiaddr"
	madns "github.com/dai/go-ipfs/gxlibs/github.com/multiformats/go-multiaddr-dns"
	"github.com/dai/go-ipfs/gxlibs/github.com/multiformats/go-multiaddr-net"
	"github.com/dai/go-ipfs/gxlibs/github.com/prometheus/client_golang/prometheus"
	"github.com/dai/go-ipfs/namesys"
	"github.com/dai/go-ipfs/plugin/loader"
	fsrepo "github.com/dai/go-ipfs/repo/fsrepo"
	"github.com/dai/types"
	"github.com/davecgh/go-spew/spew"
)

var dnsResolver = madns.DefaultResolver
var daios *dcore.Daios
var i int

const (
	nBitsForKeypairDefault = 2048
	dhtKey                 = "daios"
	blockTopic             = "daiosBlock"
	txTopic                = "daiosTx"
	peerTopic              = "daiosPeer"
	mineTopic              = "daiosMine"
)

func main() {
	ctx := context.Background()
	var err error

	intrh, ctx := setupInterruptHandler(ctx)
	defer intrh.Close()

	os.Args[0] = "ipfs"
	repoPath := os.Args[1]

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		fmt.Println(os.Args[1], "does not exist")

		err = LoadPlug(repoPath)
		if err != nil {
			panic(err)
		}

		if !fsrepo.IsInitialized(repoPath) {
			err := initWithDefaults(os.Stdout, repoPath, "")
			if err != nil {
				panic(err)
			}
		}

		/*
			err = makeSwarmKey(repoPath)
			if err != nil {
				log.Fatalf("Failed makeSwarmKey: %v", err)
			}
		*/

	} else {
		fmt.Println(os.Args[1])
		err = LoadPlug(repoPath)
		if err != nil {
			panic(err)
		}

	}

	cfg, err := loadConfig(repoPath)

	if err != nil {
		log.Fatalf("Failed loadConf: &v", err)
	}

	cctx := &types.Context{
		Path:   repoPath,
		Config: cfg,
		Node:   nil,
		Ctx:    ctx,
	}

	daemonLocked, err := fsrepo.LockedByOtherProcess(cctx.Path)
	if err != nil {
		panic(err)
	}

	if daemonLocked {
		fmt.Println("ipfs daemon is running")
		panic("ipfs daemon is running")
	}

	err = run(*cctx)
	if err != nil {
		panic(err)
	}
}

func LoadPlug(repoPath string) error {

	pluginpath := filepath.Join(repoPath, "plugins")

	_, err := checkPermissions(repoPath)
	if err != nil {
		log.Fatalf("Failed checkPermissions: %v", err)
		return err
	}

	if _, err := loader.LoadPlugins(pluginpath); err != nil {
		log.Fatalf("Failed LoadPlugins: %v", err)
		return err
	}

	return nil
}

func makeSwarmKey(repoPath string) error {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatalln("While trying to read random source:", err)
	}

	data := []byte("/key/swarm/psk/1.0.0/\n/base16/\n" + hex.EncodeToString(key))

	keypath := path.Join(repoPath, "swarm.key")

	err = ioutil.WriteFile(keypath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
func checkPermissions(path string) (bool, error) {
	_, err := os.Open(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if os.IsPermission(err) {
		return false, fmt.Errorf("error opening repository at %s: permission denied", path)
	}

	return true, nil
}

func resolveAddr(ctx context.Context, addr ma.Multiaddr) (ma.Multiaddr, error) {
	ctx, cancelFunc := context.WithTimeout(ctx, 10*time.Second)
	defer cancelFunc()

	addrs, err := dnsResolver.Resolve(ctx, addr)
	if err != nil {
		return nil, err
	}

	if len(addrs) == 0 {
		return nil, errors.New("non-resolvable API endpoint")
	}

	return addrs[0], nil
}

func loadConfig(path string) (*config.Config, error) {
	return fsrepo.ConfigAt(path)
}

type IntrHandler struct {
	sig chan os.Signal
	wg  sync.WaitGroup
}

func NewIntrHandler() *IntrHandler {
	ih := &IntrHandler{}
	ih.sig = make(chan os.Signal, 1)
	return ih
}

func (ih *IntrHandler) Close() error {
	close(ih.sig)
	ih.wg.Wait()
	return nil
}

func (ih *IntrHandler) Handle(handler func(count int, ih *IntrHandler), sigs ...os.Signal) {
	signal.Notify(ih.sig, sigs...)
	ih.wg.Add(1)
	go func() {
		defer ih.wg.Done()
		count := 0
		for range ih.sig {
			count++
			handler(count, ih)
		}
		signal.Stop(ih.sig)
	}()
}

func setupInterruptHandler(ctx context.Context) (io.Closer, context.Context) {
	intrh := NewIntrHandler()
	ctx, cancelFunc := context.WithCancel(ctx)

	handlerFunc := func(count int, ih *IntrHandler) {
		switch count {
		case 1:
			fmt.Println()

			ih.wg.Add(1)
			go func() {
				defer ih.wg.Done()
				cancelFunc()
			}()

		default:
			fmt.Println("Received another interrupt")
			os.Exit(-1)
		}
	}

	intrh.Handle(handlerFunc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	return intrh, ctx
}

func run(ctx types.Context) error {
	go func() {
		<-ctx.Ctx.Done()
		fmt.Println("Received interrupt signal")
	}()

	repo, err := fsrepo.Open(ctx.Path)
	if err != nil {
		log.Fatalf("Failed to fsrepo.Open:: %v", err)
	}

	ncfg := &core.BuildCfg{
		Repo:                        repo,
		Permanent:                   true,
		Online:                      !false,
		DisableEncryptedConnections: false,
		ExtraOpts: map[string]bool{
			"pubsub": true,
			"ipnsps": false,
			"mplex":  true,
		},
	}
	ncfg.Routing = core.DHTOption

	ctx.Node, err = core.NewNode(ctx.Ctx, ncfg)

	address := types.NewAddress(com.Hash(string(ctx.Node.PeerHost.ID().Pretty())))
	fmt.Println("Address :", types.NewAddress(com.Hash(string(ctx.Node.PeerHost.ID().Pretty()))))

	daios = dcore.New(address)

	if err != nil {
		log.Fatalf("error from node construction: ", err)
		return err
	}
	ctx.Node.SetLocal(false)

	if ctx.Node.PNetFingerprint != nil {
		fmt.Printf("Swarm key fingerprint: %x\n", ctx.Node.PNetFingerprint)
	}

	defer func() {
		ctx.Node.Close()

		select {
		case <-ctx.Ctx.Done():
			log.Println("shutdown")
		default:
		}
	}()
	ctx.Node.PeerHost.SetStreamHandler("/ipfs/id/1.0.0", handleStream)

	apiErrc, err := serveHTTPApi(ctx)
	if err != nil {
		return err
	}

	subsErrc, err := subscribe(ctx.Node, address)
	if err != nil {
		return err
	}

	gcErrc, err := runGC(ctx)
	if err != nil {
		return err
	}

	gwErrc, err := serveHTTPGateway(ctx)
	if err != nil {
		return err
	}

	prometheus.MustRegister(&corehttp.IpfsNodeCollector{Node: ctx.Node})

	if err != nil {
		log.Fatalf("Failed to start IPFS node: %v", err)
	}

	for err := range merge(apiErrc, gcErrc, gwErrc, subsErrc) {
		if err != nil {
			return err
		}
	}

	return nil
}

func subscribe(node *core.IpfsNode, address types.Address) (<-chan error, error) {
	errc := make(chan error)

	jsonBytes, err := json.Marshal(address)
	if err != nil {
		return nil, err
	}
	if err := node.PubSub.Publish(peerTopic, jsonBytes); err != nil {
		return nil, err
	}

	subBlock, err := node.PubSub.Subscribe(blockTopic)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			msg, err := subBlock.Next(node.Context())
			if err != nil {
				errc <- err
			}

			var b types.Block

			daios.TxPool().Remove()
			json.Unmarshal(msg.GetData(), &b)
			daios.BlockChain().AddBlock(&b)
			spew.Dump(b.States)
			sdb := *types.CS.DB()
			sdb.Sync(b.States)

		}
	}()

	go func() {
		for {

			broadCastBlocks := dcore.BroadCastBlocks()
			b := <-*broadCastBlocks
			i = 0

			if err := node.PubSub.Publish(blockTopic, b.MarshalJSON()); err != nil {
				errc <- err
			}

		}
	}()

	subTx, err := node.PubSub.Subscribe(txTopic)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			msg, err := subTx.Next(node.Context())
			if err != nil {
				errc <- err
			}
			var tx types.Transaction
			json.Unmarshal(msg.GetData(), &tx)

			txp := *daios.TxPool()
			txp.Enqueue(&tx)

		}
	}()

	subMine, err := node.PubSub.Subscribe(mineTopic)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			msg, err := subMine.Next(node.Context())
			if err != nil {
				errc <- err
			}

			var addr types.Address

			json.Unmarshal(msg.GetData(), &addr)

			if bytes.Equal(address[:], addr[:]) {
				dcore.Mine(daios, address)
			}

		}
	}()

	return nil, nil
}

func runGC(ctx types.Context) (<-chan error, error) {
	enableGC := true
	if !enableGC {
		return nil, nil
	}

	errc := make(chan error)
	go func() {
		errc <- corerepo.PeriodicGC(ctx.Ctx, ctx.Node)
		close(errc)
	}()
	return errc, nil
}

func serveHTTPApi(ctx types.Context) (<-chan error, error) {
	apiAddrs := make([]string, 0, 2)

	apiAddrs = ctx.Config.Addresses.API

	listeners := make([]manet.Listener, 0, len(apiAddrs))
	for _, addr := range apiAddrs {
		apiMaddr, err := ma.NewMultiaddr(addr)
		if err != nil {
			return nil, fmt.Errorf("serveHTTPApi: invalid API address: %q (err: %s)", apiMaddr, err)
		}
		apiLis, err := manet.Listen(apiMaddr)

		if err != nil {
			return nil, fmt.Errorf("serveHTTPApi: manet.Listen(%s) failed: %s", apiMaddr, err)
		}
		apiMaddr = apiLis.Multiaddr()

		fmt.Printf("API server listening on %s\n", apiMaddr)

		listeners = append(listeners, apiLis)
	}

	gatewayOpt := corehttp.GatewayOption(true, "/ipfs", "/ipns")
	dirb := uio.NewDirectory(ctx.Node.DAG)

	var opts = []corehttp.ServeOption{
		corehttp.MetricsCollectionOption("api"),
		corehttp.CheckVersionOption(),
		gatewayOpt,
		uploadMux("/upload", dirb),
		downloadMux("/download"),
		listMux("/list"),
		blockMux("/block"),
		statsMux("/stats"),
	}

	APIPath := "/api/v0"
	originEnvKey := "API_ORIGIN"
	var defaultLocalhostOrigins = []string{
		"http://127.0.0.1:<port>",
		"https://127.0.0.1:<port>",
		"http://localhost:<port>",
		"https://localhost:<port>",
	}

	cfg := cmdsHttp.NewServerConfig()

	cfg.APIPath = APIPath

	rcfg, err := ctx.Node.Repo.Config()
	if err != nil {
		return nil, err
	}

	if acao := rcfg.API.HTTPHeaders[cmdsHttp.ACAOrigin]; acao != nil {
		cfg.SetAllowedOrigins(acao...)
	}
	if acam := rcfg.API.HTTPHeaders[cmdsHttp.ACAMethods]; acam != nil {
		cfg.SetAllowedMethods(acam...)
	}
	if acac := rcfg.API.HTTPHeaders[cmdsHttp.ACACredentials]; acac != nil {
		for _, v := range acac {
			cfg.SetAllowCredentials(strings.ToLower(v) == "true")
		}
	}

	cfg.Headers = make(map[string][]string, len(rcfg.API.HTTPHeaders))

	for h, v := range rcfg.API.HTTPHeaders {
		cfg.Headers[h] = v
	}
	cfg.Headers["Server"] = []string{"go-ipfs/" + version.CurrentVersionNumber}

	origin := os.Getenv(originEnvKey)
	if origin != "" {
		cfg.AppendAllowedOrigins(origin)
	}

	if len(cfg.AllowedOrigins()) == 0 {
		cfg.SetAllowedOrigins(defaultLocalhostOrigins...)
	}

	if len(cfg.AllowedMethods()) == 0 {
		cfg.SetAllowedMethods("GET", "POST", "PUT")
	}

	port := ""
	if tcpaddr, ok := listeners[0].Addr().(*net.TCPAddr); ok {
		port = strconv.Itoa(tcpaddr.Port)
	} else if udpaddr, ok := listeners[0].Addr().(*net.UDPAddr); ok {
		port = strconv.Itoa(udpaddr.Port)
	}

	oldOrigins := cfg.AllowedOrigins()
	newOrigins := make([]string, len(oldOrigins))
	for i, o := range oldOrigins {

		if port != "" {
			o = strings.Replace(o, "<port>", port, -1)
		}
		newOrigins[i] = o
	}
	cfg.SetAllowedOrigins(newOrigins...)

	if err := ctx.Node.Repo.SetAPIAddr(listeners[0].Multiaddr()); err != nil {
		return nil, fmt.Errorf("serveHTTPApi: SetAPIAddr() failed: %s", err)
	}

	errc := make(chan error)
	var wg sync.WaitGroup
	for _, apiLis := range listeners {
		wg.Add(1)
		go func(lis manet.Listener) {
			defer wg.Done()
			errc <- corehttp.Serve(ctx.Node, manet.NetListener(lis), opts...)

		}(apiLis)

	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	return errc, nil
}

type info struct {
	Stat  corerepo.SizeStat `json:"Repo"`
	Stats metrics.Stats     `json:"Metrics"`
}

func statsMux(path string) corehttp.ServeOption {
	return func(node *core.IpfsNode, _ net.Listener, mux *http.ServeMux) (*http.ServeMux, error) {
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				return
			}
			//디스크 사용량
			stat, err := corerepo.RepoStat(node.Context(), node)
			if err != nil {
				fmt.Println("Failed Repostat")
			}
			fmt.Printf("current : %d byte\n", stat.RepoSize)
			fmt.Printf("max : %d byte\n", stat.StorageMax)

			totals := node.Reporter.GetBandwidthTotals()
			fmt.Printf("TotalIn: %s\n", humanize.Bytes(uint64(totals.TotalIn)))
			fmt.Printf("TotalOut: %s\n", humanize.Bytes(uint64(totals.TotalOut)))
			fmt.Printf("RateIn: %s/s\n", humanize.Bytes(uint64(totals.RateIn)))
			fmt.Printf("RateOut: %s/s\n", humanize.Bytes(uint64(totals.RateOut)))
			info := &info{
				Stat:  stat.SizeStat,
				Stats: totals,
			}
			data, err := json.Marshal(info)
			if err != nil {
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write(data)
		})

		return mux, nil
	}
}

func blockMux(path string) corehttp.ServeOption {
	return func(node *core.IpfsNode, _ net.Listener, mux *http.ServeMux) (*http.ServeMux, error) {
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				return
			}
			jsonChain, err := json.Marshal(daios.BlockChain())
			if err != nil {
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write(jsonChain)

		})
		return mux, nil
	}
}

func uploadMux(path string, dirb uio.Directory) corehttp.ServeOption {
	return func(node *core.IpfsNode, _ net.Listener, mux *http.ServeMux) (*http.ServeMux, error) {
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {

			if r.Method != "POST" {
				return
			}

			var fileHeader *multipart.FileHeader
			var e error
			file, fileHeader, e := r.FormFile("file")
			if e != nil {
				fmt.Errorf("Failed FormFile:", e)
				return
			}

			s, err := coreunix.Add(node, file)
			if err != nil {
				fmt.Errorf("Failed coreunix.Add", err)
				return
			}

			c, err := cid.Decode(s)
			if err != nil {
				fmt.Errorf("Failed cid.Decode", err)
				return
			}

			n, err := node.DAG.Get(node.Context(), c)
			if err != nil {
				fmt.Errorf("Failed node.DAG.Get", err)
				return
			}

			if err := dirb.AddChild(node.Context(), fileHeader.Filename, n); err != nil {
				fmt.Errorf("Failed dirb.AddChild", err)
				return
			}
			dir, err := dirb.GetNode()

			if err := node.Pinning.Pin(node.Context(), dir, true); err != nil {
				fmt.Errorf("Failed dirb.GetNode", err)
				return
			}
			if err := node.Pinning.Flush(); err != nil {
				fmt.Errorf("Failed dirb.GetNode", err)
				return
			}

			data, err := json.Marshal(dir.Cid())
			fmt.Println("dir:", dir.Cid())

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write(data)

			var mutex sync.Mutex

			payload := s + "/" + fileHeader.Filename

			fmt.Println("payload: ", payload)

			mutex.Lock()
			d := types.NewTransaction(types.NewAddress(""), types.NewAddress(""), 0, payload)

			d.Nonce = i
			d.Data.Hash = d.Hash()
			mutex.Unlock()

			if err := node.PubSub.Publish(txTopic, d.MarshalJSON()); err != nil {
				return
			}
			i++

		})
		return mux, nil
	}
}

func downloadMux(key string) corehttp.ServeOption {
	return func(node *core.IpfsNode, _ net.Listener, mux *http.ServeMux) (*http.ServeMux, error) {
		mux.HandleFunc(key, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				return
			}

			k := r.FormValue("key")
			n := r.FormValue("name")

			rd, err := coreunix.Cat(node.Context(), node, k)

			if err != nil {
				log.Fatalf("Failed coreunix.Cat: %v", err)
				return
			}
			w.Header().Set("Content-Disposition", "attachment; filename="+n)

			io.Copy(w, rd)

		})
		return mux, nil
	}
}

func listMux(key string) corehttp.ServeOption {
	return func(node *core.IpfsNode, _ net.Listener, mux *http.ServeMux) (*http.ServeMux, error) {
		mux.HandleFunc(key, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				return
			}
			api := coreapi.NewCoreAPI(node)

			var paths []string

			paths = append(paths, r.FormValue("key"))

			var dagnodes []ipld.Node
			for _, fpath := range paths {
				p, err := iface.ParsePath(fpath)
				fmt.Println(p)
				if err != nil {
					fmt.Errorf("Failed: iface.ParsePath", err)
					return
				}

				dagnode, err := api.ResolveNode(node.Context(), p)
				if err != nil {
					fmt.Errorf("Failed: api.ResolveNode", err)
					return
				}
				dagnodes = append(dagnodes, dagnode)
			}

			output := make([]ls.LsObject, len(paths))

			ng := merkledag.NewSession(node.Context(), node.DAG)

			ro := merkledag.NewReadOnlyDagService(ng)

			for i, dagnode := range dagnodes {
				dir, err := uio.NewDirectoryFromNode(ro, dagnode)

				if err != nil && err != uio.ErrNotADir {
					fmt.Errorf("Failed: uio", err)
					return
				}

				var links []*ipld.Link
				if dir == nil {
					links = dagnode.Links()
				} else {
					links, err = dir.Links(node.Context())
					if err != nil {
						fmt.Errorf("Failed: Links", err)
						return
					}
				}

				output[i] = ls.LsObject{
					Hash:  paths[i],
					Links: make([]ls.LsLink, len(links)),
				}

				for j, link := range links {
					t := unixfspb.Data_DataType(-1)

					switch link.Cid.Type() {
					case cid.Raw:

						t = unixfs.TFile
					case cid.DagProtobuf:
						linkNode, err := link.GetNode(node.Context(), node.DAG)
						if err == ipld.ErrNotFound {

							linkNode = nil
						} else if err != nil {
							fmt.Errorf("Failed: GetNode", err)
							return

						}

						if pn, ok := linkNode.(*merkledag.ProtoNode); ok {
							d, err := unixfs.FSNodeFromBytes(pn.Data())
							if err != nil {
								fmt.Errorf("Failed: FSNodeFromBytes", err)
								return
							}
							t = d.Type()
						}
					}
					output[i].Links[j] = ls.LsLink{
						Name: link.Name,
						Hash: link.Cid.String(),
						Size: link.Size,
						Type: t,
					}

					fmt.Println(output[i].Links[j])
				}
			}

			data, err := json.Marshal(output)
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write(data)

			if err != nil {
				panic(err)
			}

		})
		return mux, nil
	}
}

func merge(cs ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	out := make(chan error)

	output := func(c <-chan error) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	for _, c := range cs {
		if c != nil {
			wg.Add(1)
			go output(c)
		}
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func serveHTTPGateway(ctx types.Context) (<-chan error, error) {

	gatewayAddrs := make([]string, 0, 2)
	gatewayAddrs = ctx.Config.Addresses.Gateway
	writable := ctx.Config.Gateway.Writable

	listeners := make([]manet.Listener, 0, len(gatewayAddrs))
	for _, addr := range gatewayAddrs {
		gatewayMaddr, err := ma.NewMultiaddr(addr)
		if err != nil {
			return nil, fmt.Errorf("serveHTTPGateway: invalid gateway address: %q (err: %s)", addr, err)
		}

		gwLis, err := manet.Listen(gatewayMaddr)
		if err != nil {
			return nil, fmt.Errorf("serveHTTPGateway: manet.Listen(%s) failed: %s", gatewayMaddr, err)
		}

		gatewayMaddr = gwLis.Multiaddr()

		if writable {
			fmt.Printf("Gateway (writable) server listening on %s\n", gatewayMaddr)
		} else {
			fmt.Printf("Gateway (readonly) server listening on %s\n", gatewayMaddr)
		}
		listeners = append(listeners, gwLis)
	}

	var opts = []corehttp.ServeOption{
		corehttp.MetricsCollectionOption("gateway"),
		corehttp.IPNSHostnameOption(),
		corehttp.GatewayOption(false, "/ipfs", "/ipns"),
		corehttp.CheckVersionOption(),
	}

	errc := make(chan error)
	var wg sync.WaitGroup
	for _, lis := range listeners {
		wg.Add(1)
		go func(lis manet.Listener) {
			defer wg.Done()
			errc <- corehttp.Serve(ctx.Node, manet.NetListener(lis), opts...)
		}(lis)
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	return errc, nil
}

func initWithDefaults(out io.Writer, repoRoot string, profile string) error {
	var profiles []string
	if profile != "" {
		profiles = strings.Split(profile, ",")
	}

	return doInit(out, repoRoot, false, nBitsForKeypairDefault, profiles, nil)
}

func doInit(out io.Writer, repoRoot string, empty bool, nBitsForKeypair int, confProfiles []string, conf *config.Config) error {
	if _, err := fmt.Fprintf(out, "path: %s\n", repoRoot); err != nil {
		return err
	}

	if fsrepo.IsInitialized(repoRoot) {
		return errors.New("Failed IsInitialized")
	}

	if conf == nil {
		var err error
		conf, err = config.Init(out, nBitsForKeypair)
		if err != nil {
			return err
		}
	}

	for _, profile := range confProfiles {
		transformer, ok := config.Profiles[profile]
		if !ok {
			return fmt.Errorf("invalid configuration profile: %s", profile)
		}

		if err := transformer.Transform(conf); err != nil {
			return err
		}
	}

	if err := fsrepo.Init(repoRoot, conf); err != nil {
		return err
	}

	if !empty {
		if err := addDefaultAssets(out, repoRoot); err != nil {
			return err
		}
	}

	return initializeIpnsKeyspace(repoRoot)
}

func addDefaultAssets(out io.Writer, repoRoot string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r, err := fsrepo.Open(repoRoot)
	if err != nil {
		return err
	}

	nd, err := core.NewNode(ctx, &core.BuildCfg{Repo: r})

	if err != nil {
		return err
	}

	defer nd.Close()

	return err
}

func initializeIpnsKeyspace(repoRoot string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r, err := fsrepo.Open(repoRoot)
	if err != nil {
		return err
	}

	nd, err := core.NewNode(ctx, &core.BuildCfg{Repo: r})
	if err != nil {
		return err
	}
	defer nd.Close()

	err = nd.SetupOfflineRouting()
	if err != nil {
		return err
	}

	return namesys.InitializeKeyspace(ctx, nd.Namesys, nd.Pinning, nd.PrivateKey)
}

func handleStream(s inet.Stream) {
	fmt.Println(s.Conn().RemotePeer().Pretty)
	defer s.Close()
	var mutex sync.Mutex

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	mutex.Lock()
	bc := *daios.BlockChain()
	rw.WriteString(string(bc.MarshalJSON()) + "\n")
	rw.Flush()
	mutex.Unlock()
}
