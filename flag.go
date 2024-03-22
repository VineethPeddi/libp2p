package main

import (
	"flag"
	"fmt"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	maddr "github.com/multiformats/go-multiaddr"
	"strings"
)

type addrList []maddr.Multiaddr

func (al *addrList) String() string {
	strs := make([]string, len(*al))
	for i, addr := range *al {
		strs[i] = addr.String()
	}
	return strings.Join(strs, ",")
}

func (al *addrList) Set(value string) error {
	addr, err := maddr.NewMultiaddr(value)
	if err != nil {
		return err
	}
	*al = append(*al, addr)
	return nil
}

func StringsToAddrs(addrStrings []string) (maddrs []maddr.Multiaddr, err error) {
	for _, addrString := range addrStrings {
		addr, err := maddr.NewMultiaddr(addrString)
		if err != nil {
			return maddrs, err
		}
		maddrs = append(maddrs, addr)
	}
	return
}

type Config struct {
	RendezvousString string
	ProtocolID       string
	BootstrapPeers   addrList
	listenHost       string
	listenPort       int
	dType            string
}

func parseFlags() *Config {
	c := &Config{}

	flag.StringVar(&c.RendezvousString, "rendezvous", "BlockChains", "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	flag.StringVar(&c.listenHost, "host", "0.0.0.0", "The bootstrap node host listen address\n")
	flag.StringVar(&c.ProtocolID, "pid", "/chat/1.1.0", "Sets a protocol id for stream headers")
	flag.IntVar(&c.listenPort, "port", 4001, "node listen port")
	flag.StringVar(&c.dType, "dType", "mdns", "Discovery type")
	flag.Var(&c.BootstrapPeers, "peer", "Adds a peer multiaddress to the bootstrap list")

	flag.Parse()
	if len(c.BootstrapPeers) == 0 {
		c.BootstrapPeers = dht.DefaultBootstrapPeers
	}

	if err := c.validateConfig(); err != nil {
		panic(err)
	}
	return c
}

func (c *Config) validateConfig() error {
	if c.dType != "mdns" && c.dType != "dht" {
		return fmt.Errorf("Invalid discovery type %v . Please use either 'mdns' or 'dht'", c.dType)
	}
	return nil
}
