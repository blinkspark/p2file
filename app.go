package p2file

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/blinkspark/p2file/config"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

type App struct {
	Host    host.Host
	Dht     *dht.IpfsDHT
	dirName string
}

type NewAppArgs struct {
	BootstrapPeers []peer.AddrInfo
}

type NewAppOpts func(arg *NewAppArgs)

var DefaultAppArgs = NewAppArgs{
	BootstrapPeers: dht.GetDefaultBootstrapPeerAddrInfos(),
}

func WithBootstrapPeers(bootstrapPeers []peer.AddrInfo) NewAppOpts {
	return func(arg *NewAppArgs) {
		arg.BootstrapPeers = bootstrapPeers
	}
}

// NewApp 创建并返回一个新的应用实例
// 参数 args 是可变参数，用于配置应用实例的选项
// 返回值：
//   - *App: 成功时返回应用实例的指针
//   - error: 如果创建过程中出现错误，则返回错误信息
func NewApp(opts ...NewAppOpts) (*App, error) {
	// 使用默认的应用参数初始化
	args := DefaultAppArgs

	// 遍历所有传入的选项，并应用到应用参数上
	for _, opt := range opts {
		opt(&args)
	}

	// 创建一个新的 libp2p 主机，启用自动NAT、NAT服务、端口映射和打洞功能
	host, err := libp2p.New(libp2p.EnableAutoNATv2(), libp2p.EnableNATService(), libp2p.NATPortMap(), libp2p.EnableHolePunching())
	if err != nil {
		return nil, err
	}

	// 使用主机和引导节点创建一个新的DHT节点
	dntNode, err := dht.New(context.Background(), host, dht.BootstrapPeers(args.BootstrapPeers...))
	if err != nil {
		return nil, err
	}

	// 启动DHT节点的引导过程
	err = dntNode.Bootstrap(context.Background())
	if err != nil {
		return nil, err
	}

	// 创建应用实例，包含libp2p主机和DHT节点
	app := &App{
		Host: host,
		Dht:  dntNode,
	}

	return app, nil
}

// WaitBootstrap 等待应用完成引导过程，确保网络连接达到最小要求。
// 参数:
//   - timeoutSec: 超时时间（秒），如果小于等于0则默认使用10秒。
//
// 返回值:
//   - error: 如果在超时时间内未能建立足够的网络连接，则返回错误；否则返回nil。
func (app *App) WaitBootstrap(timeoutSec int) error {
	if timeoutSec <= 0 {
		timeoutSec = 10
	}

	timeout := time.After(time.Duration(timeoutSec) * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("bootstrap timeout after %d seconds", timeoutSec)
		case <-ticker.C:
			if len(app.Host.Network().Peers()) > 5 {
				return nil
			}
		}
	}
}

func (app *App) ListDir() ([]string, error) {
	id, err := peer.Decode(*config.Channel)
	if err != nil {
		return nil, err
	}
	pi, err := app.Dht.FindPeer(context.Background(), id)
	if err != nil {
		return nil, err
	}
	log.Println(pi)
	err = app.Host.Connect(context.Background(), pi)
	if err != nil {
		return nil, err
	}

	topic := "/p2file/" + id.String()
	log.Println("connecting to " + topic)
	stream, err := app.Host.NewStream(context.Background(), id, protocol.ID(topic))
	if err != nil {
		return nil, err
	}
	defer stream.Close()
	log.Println("connected to " + topic)
	reader := bufio.NewReader(stream)
	writer := bufio.NewWriter(stream)

	payload := Payload{
		Type: PL_LS,
	}

	log.Printf("sending payload %+v\n", payload)
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	_, err = writer.Write(raw)
	if err != nil {
		return nil, err
	}
	writer.WriteByte('\n')
	writer.Flush()

	raw, err = reader.ReadBytes('\n')
	err = json.Unmarshal(raw, &payload)
	if err != nil {
		return nil, err
	}
	log.Printf("received payload %+v\n", payload)

	return payload.DirList, nil
}

func (app *App) Serve(dirName string) error {
	app.dirName = dirName
	// listen terminal signal to close the host
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	topic := "/p2file/" + app.Host.ID().String()
	app.Host.SetStreamHandler(protocol.ID(topic), app.handleServe)

	<-sigchan
	return nil
}

func (app *App) handleServe(stream network.Stream) {
	log.Println("handle serve")
	defer stream.Close()
	reader := bufio.NewReader(stream)
	writer := bufio.NewWriter(stream)
	for {
		raw, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Println("read error:", err)
			return
		}

		payload := Payload{}
		err = json.Unmarshal(raw, &payload)
		if err != nil {
			log.Println("json error:", err)
			return
		}

		switch payload.Type {
		case PL_LS:
			files, err := os.ReadDir(app.dirName)
			if err != nil {
				log.Println("read dir error:", err)
				return
			}
			resPayload := Payload{
				Type: PL_LS_RES,
			}
			for _, file := range files {
				if file.IsDir() {
					continue
				}
				name := file.Name()
				resPayload.DirList = append(resPayload.DirList, name)
			}
			res, err := json.Marshal(resPayload)
			if err != nil {
				log.Println("json error:", err)
				return
			}
			writer.Write(res)
			writer.WriteByte('\n')
			writer.Flush()
			log.Println("send done")
		case PL_GET:
			fileName := payload.TargetFile
			// file, err := os.Open(app.dirName + "/" + fileName)
			file, err := os.Open(path.Join(app.dirName, fileName))
			if err != nil {
				log.Println("open file error:", err)
				return
			}
			defer file.Close()

			buf := make([]byte, 1024)
			for {
				resPayload := Payload{
					Type: PL_GET_RES,
				}
				n, err := file.Read(buf)
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Println("read file error:", err)
					return
				}
				resPayload.Data = buf[:n]
				res, err := json.Marshal(resPayload)
				if err != nil {
					log.Println("json error:", err)
					return
				}
				writer.Write(res)
				writer.WriteByte('\n')
			}
			payload := Payload{
				Type: PL_GET_RES_DONE,
			}
			res, err := json.Marshal(payload)
			if err != nil {
				log.Println("json error:", err)
				return
			}
			writer.Write(res)
			writer.WriteByte('\n')
			writer.Flush()
		}
	}
}

func (app *App) Close() error {
	return app.Host.Close()
}
