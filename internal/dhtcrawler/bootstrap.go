package dhtcrawler

import (
	"context"
	"net"
	"net/netip"
	"regexp"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
)

func (c *crawler) reseedBootstrapNodes(ctx context.Context) {
	interval := time.Duration(0)

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			for _, strAddr := range c.bootstrapNodes {
				matchIpNumber, _ := regexp.MatchString(`\d+\.\d+\.\d+\.\d+:\d+`, strAddr)
				var addrPort netip.AddrPort
				if matchIpNumber {
					var err error
					addrPort, err = netip.ParseAddrPort(strAddr)
					if err != nil {
						c.logger.Warnf("invalid IP or port number: %s", strAddr)
						continue
					}
				} else {
					addr, err := net.ResolveUDPAddr("udp", strAddr)
					if err != nil {
						c.logger.Warnf("failed to resolve bootstrap node address: %s", err)
						continue
					}
					addrPort = addr.AddrPort()
				}

				select {
				case <-ctx.Done():
					return
				case c.nodesForPing.In() <- ktable.NewNode(ktable.ID{}, addrPort):
					continue
				}
			}
		}

		interval = c.reseedBootstrapNodesInterval
	}
}
