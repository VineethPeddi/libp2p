This is a peer to peer file sharing app using streams.

Build the project using the following command

`go build -o p2p`

Execute the build using the following command:

1. For MDNS discovery

        In terminal 1: `./p2p -port 4001 -dType "mdns" `
        
        In terminal 2: `./p2p -port 4002 -dType "mdns" `
        
        Now we are running two peers in different port successfully. These two peers will be connected and will be ready to share files between each other.
2. For DHT discovery
        
        In terminal 1: `./p2p -port 4001 -dType "dht" `
        
        In terminal 2: `./p2p -port 4002 -dType "dht" `

Note: `If "dType" flag is not provided then by default "mdns" discovery is executed"` 


For every new peer connection, a stream will be created and files can be shared over this stream.
We can upload any text file from a peer and other peer can download the file. However we can customize to upload files of any type.

The file will be downloaded by the name `output.txt`.

