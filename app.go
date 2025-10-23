package p2file

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type App struct {
	Host host.Host
	Dht  *dht.IpfsDHT
}

type NewAppArgs struct {
	BootstrapPeers []peer.AddrInfo
}

func NewApp(args *NewAppArgs) (*App, error) {
	if args == nil {
		args = &NewAppArgs{}
	}

	// create a new libp2p host
	host, err := libp2p.New(libp2p.EnableAutoNATv2(), libp2p.EnableNATService(), libp2p.NATPortMap(), libp2p.EnableHolePunching())
	if err != nil {
		return nil, err
	}

	if args.BootstrapPeers == nil {
		args.BootstrapPeers = dht.GetDefaultBootstrapPeerAddrInfos()
	}

	dntNode, err := dht.New(context.Background(), host, dht.BootstrapPeers(args.BootstrapPeers...))
	if err != nil {
		return nil, err
	}

	err = dntNode.Bootstrap(context.Background())
	if err != nil {
		return nil, err
	}

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

func (app *App) Serve() error {
	// listen terminal signal to close the host
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	<-sigchan
	return nil
}

func (app *App) Close() error {
	return app.Host.Close()
}
