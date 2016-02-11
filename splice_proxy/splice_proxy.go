package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"syscall"
)

const (
	// These constants are pulled from kernel source.
	SPLICE_F_MOVE     = 0x01
	SPLICE_F_MORE     = 0x04
	SPLICE_F_NONBLOCK = 002
)

var (
	inPort     = flag.Int("in-port", 8081, "input port number.")
	outPort    = flag.Int("out-port", 3306, "output port number.")
	bufferSize = flag.Int("buff-size", 16*1024, "splice buffer size.")
)

func main() {
	flag.Parse()
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", *inPort))
	if err != nil {
		log.Fatal(err)
	}
	// TODO: Close fds
	for {
		in, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		out, err := net.Dial("tcp", fmt.Sprintf("localhost:%v", *outPort))
		if err != nil {
			log.Fatal(err)
		}
		fromBackend := NewTCPCopier()
		go fromBackend.CopyBytes("from backend", in.(*net.TCPConn), out.(*net.TCPConn), nil)
		toBackend := NewTCPCopier()
		go toBackend.CopyBytes("to backend", out.(*net.TCPConn), in.(*net.TCPConn), nil)
	}
}

type TCPCopier interface {
	CopyBytes(direction string, in, out *net.TCPConn, wg *sync.WaitGroup)
}

func NewTCPCopier() TCPCopier {
	return &Splicer{tcpFD: make([]int, 2), pipe: make([]int, 2), bufferSize: *bufferSize}
}

type Splicer struct {
	tcpFD      []int
	pipe       []int
	bufferSize int
}

func (s *Splicer) shutdown() {
	for _, f := range append(s.pipe, s.tcpFD...) {
		syscall.Close(f)
	}
}

func (s *Splicer) CopyBytes(direction string, in, out *net.TCPConn, wg *sync.WaitGroup) {
	// Signal to the caller that we're done
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()
	// Call close on the required fds.
	defer s.shutdown()
	// Close the tcp connection given to us.
	defer func() {
		// TODO: Double close problems?
		out.Close()
		in.Close()
	}()

	// Invoking File() on the connection will duplicate the socket fd, setting
	// the original (socket fd) to block.
	inFile, err := in.File()
	if err != nil {
		log.Fatal(err)
	}
	s.tcpFD[0] = int(inFile.Fd())
	outFile, err := out.File()
	if err != nil {
		log.Fatal(err)
	}
	s.tcpFD[1] = int(outFile.Fd())

	if err := syscall.Pipe(s.pipe); err != nil {
		log.Fatal(err)
	}
	log.Printf("Splicing %s: %s -> %s, pipe %+v, tcpfds %+v", direction, in.RemoteAddr(), out.RemoteAddr(), s.pipe, s.tcpFD)
	var netBytes int64
	for {
		// SPLICE_F_NONBLOCK: don't block if TCP buffer is empty
		// SPLICE_F_MOVE: directly move pages into splice buffer in kernel memory
		// SPLICE_F_MORE: just makes mysql connections slow
		inBytes, err := syscall.Splice(s.tcpFD[0], nil, s.pipe[1], nil, s.bufferSize, SPLICE_F_MOVE|SPLICE_F_NONBLOCK)
		if inBytes == 0 {
			break
		}
		if err != nil {
			if err == syscall.EAGAIN || err == syscall.EINTR {
				log.Printf("Retrying error: %v", err)
				continue
			}
			log.Fatalf("error splicing into pipe: %v", err)
		}
		outBytes, err := syscall.Splice(s.pipe[0], nil, s.tcpFD[1], nil, s.bufferSize, SPLICE_F_MOVE|SPLICE_F_NONBLOCK)
		if err != nil {
			if err == syscall.EAGAIN || err == syscall.EINTR {
				log.Printf("Retrying error: %v", err)
				continue
			}
			log.Fatalf("error splicing out of pipe: %+v", err)
		}
		log.Printf("Splicing, out %+v, in %+v, buffer size %+v", outBytes, inBytes, s.bufferSize)
		netBytes += outBytes
	}
	log.Printf("Splicing %d bytes %s: %s -> %s", netBytes, direction, in.RemoteAddr(), out.RemoteAddr())
}
