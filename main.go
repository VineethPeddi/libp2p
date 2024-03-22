package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/multiformats/go-multiaddr"
	"os"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
)

func handleStream(stream network.Stream) {
	fmt.Println("Got a new stream!")

	// Create a buffer stream for non-blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readData(rw)
	go writeData(rw)

	// 'stream' will stay open until you close it (or the other side closes it).
}

func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from buffer")
			panic(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">  ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}

		_, err = rw.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			fmt.Println("Error writing to buffer")
			panic(err)
		}
		err = rw.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer")
			panic(err)
		}
	}
}

func getDiscoveryType(dtype string) Discovery {
	if dtype == "mdns" {
		return &Mdns{}
	}
	return &DHT{}
}

func createNode(config *Config) host.Host {

	fmt.Printf("[*] Listening on: %s with port: %d\n", config.listenHost, config.listenPort)

	r := rand.Reader

	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", config.listenHost, config.listenPort))
	fmt.Println("source: ", sourceMultiAddr, "\n\n ")
	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.

	host, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("host ID: ", host.ID())
	fmt.Println("host address: ", host.Addrs())
	return host
}

func main() {

	config := parseFlags()
	host := createNode(config)

	// Set a function as stream handler.
	// This function is called when a peer initiates a connection and starts a stream with this peer.
	host.SetStreamHandler(protocol.ID(config.ProtocolID), handleStream)

	discovery := getDiscoveryType(config.dType)
	discovery.initDiscovery(host, config)

	select {}
}
