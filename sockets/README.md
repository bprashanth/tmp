# Problem statement

## Why raw sockets?

* We can only receive frames destined to us (unicast), to everyone (broadcast)
or those we subscribe to (multicast - video conf).
* All headers are stripped by network stack.
* Cannot modify packet headers, kernel prepends headers.

If we could receive the frames for all computers connected to our broadcast domain? (promiscuous)
If we could get all headers?
If we could inject packets with custom headers?

## Promiscous mode

"See all, hear all" wizard mode, mostly used by networking monitoring tools.
On linux we set a flag IFF_PROMISC on a device using an ioctl.
We can also use `ifconfig eth0 promisc` and press enter.

## Sniffing

Broadcast domain: ethernet domain in which all computers are connected
together and are contending with each other to send packets. Security and
DOS tools: make our own packets and sent out directly into the network.
Modify the source ip. Total network stack bypass.

Userland
Socket
TCP/UDP
IP
Protocol family ---> raw socket
NIC

## PF_PACKET

software interface that allows us to send/receive packets at l2, directly
to and from the device driver. Filtering supported on PF_PACKET interface using
BPF.

### Creation

Call socket(family, type, protocol): socket(PF_PACKET, SOCK_RAW, int protocol)
protocol: ETH_P_IP for IP networks, ETH_P_ALL for all.

### Sniffer

* Create socket()
* Set interface you want to sniff on in promiscous mode
* Bind Raw socket to interface - bind()
* Receive packets on the socket - recvfrom()
* Process received packets
* Close the raw socket

### Packet injector

* Create a raw socket - socket()
* Bind socket to the interface you want to send packets onto - bind()
* Create packet
* Send packet - sendto()
* Close the raw socket

Too much theory leads to more confusion.
Run via `gcc sniffer.c && sudo ./a.out` and packets get logged to log.txt.
