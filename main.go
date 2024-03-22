package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
	"io"
	"os"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
)

func handleStream(stream network.Stream) {
	fmt.Println("Got a new stream!")

	// Creating a buffer stream for non-blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readDataFromStream(rw)
	go writeDataToStream(rw)
}

func readDataFromStream(rw *bufio.ReadWriter) {
	for {
		file, err := os.OpenFile("example.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		writer := bufio.NewWriter(file)
		_, err = io.Copy(writer, rw.Reader)
		if err != nil {
			panic(err)
		}
	}
}

func writeDataToStream(rw *bufio.ReadWriter) {
	for {
		fmt.Print(">Please provide fileName to upload:  ")
		var fileName string
		_, err := fmt.Scan(&fileName)
		data, err := os.ReadFile(fileName)
		if err != nil {
			panic(err)
		}
		_, err = rw.Write(data)
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

	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", config.listenHost, config.listenPort))

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
