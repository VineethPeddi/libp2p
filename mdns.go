package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

// interface to be called when new  peer is found
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}

// Initialize the MDNS service
func (m *Mdns) initMDNS(peerhost host.Host, rendezvous string) (chan peer.AddrInfo, error) {
	n := &discoveryNotifee{}
	n.PeerChan = make(chan peer.AddrInfo)

	ser := mdns.NewMdnsService(peerhost, rendezvous, n)
	if err := ser.Start(); err != nil {
		return nil, err
	}
	return n.PeerChan, nil
}

type Mdns struct{}

func (m *Mdns) initDiscovery(host host.Host, config *Config) error {
	ctx := context.Background()
	peerChan, err := m.initMDNS(host, config.RendezvousString)
	if err != nil {
		return err
	}
	for peer := range peerChan { // allows multiple peers to join
		fmt.Printf("Received a peer %+v: \n\n", peer)
		if peer.ID > host.ID() {
			// if other end peer id greater than us, don't connect to it, just wait for it to connect us
			fmt.Println("Found peer:", peer, " id is greater than us, wait for it to connect to us")
			continue
		}
		fmt.Println("Found peer:", peer, ", connecting")

		if err := host.Connect(ctx, peer); err != nil {
			fmt.Println("Connection failed:", err)
			continue
		}

		// open a stream, this stream will be handled by handleStream other end
		stream, err := host.NewStream(ctx, peer.ID, protocol.ID(config.ProtocolID))

		if err != nil {
			fmt.Println("Stream open failed", err)
		} else {
			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

			go writeDataToStream(rw)
			go readDataFromStream(rw)
			fmt.Println("Connected to:", peer)
		}
	}
	return nil
}
