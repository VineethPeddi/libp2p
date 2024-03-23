This is a peer to peer file sharing app using streams built using web3 and libp2p.

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

We can pass different flags like the following:
1. -port => This is used when we want to run the application on a specific port
2. -rendezvous => This is used to provide the string through which we can find the peers
3. -host => This is the listen address 
4. -pid => This is the protocol id at which a peer can connect with others
5. -dType => This is used to specify which discovery mechanism we would like to implement
6. -peer => This is the bootstrap nodes address

