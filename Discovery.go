package main

import "github.com/libp2p/go-libp2p/core/host"

type Discovery interface {
	initDiscovery(host host.Host, config *Config) error
}
