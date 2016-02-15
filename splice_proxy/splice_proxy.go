package main

import (
	"flag"
	"fmt"
	"github.com/golang/go/src/io"
	"github.com/golang/go/src/net"
	"log"
	"strings"
	"sync"
	"syscall"
)

const (
	// These constants are pulled from kernel source.
	SPLICE_F_MOVE     = 0x01
	SPLICE_F_MORE     = 0x04
	SPLICE_F_NONBLOCK = 002
	SHUT_RDWR         = 2
	SHUT_WR           = 1
	SHUT_RD           = 0
)

var (
	inPort     = flag.Int("in-port", 8081, "input port number.")
	outPort    = flag.Int("out-port", 3306, "output port number.")
	bufferSize = flag.Int("buff-size", 16*1024, "splice buffer size.")
	splice     = flag.Bool("splice", true, "if true, use splice, otherwise use io.Copy.")
)

func main() {
	flag.Parse()
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", *inPort))
	if err != nil {
		log.Fatal(err)
	}
	// TODO: Close fds
	for {
		log.Printf("Accepting input connection")
		in, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		out, err := net.Dial("tcp", fmt.Sprintf("localhost:%v", *outPort))
		if err != nil {
			log.Fatal(err)
		}
		NewTCPCopier(*splice).CopyBytes(in.(*net.TCPConn), out.(*net.TCPConn))
	}
}

type TCPCopier interface {
	CopyBytes(in, out *net.TCPConn)
}

func NewTCPCopier(splice bool) TCPCopier {
	if splice {
		return &Splicer{tcpFD: make([]int, 2), pipe: make([]int, 2), bufferSize: *bufferSize}
	}
	return &SimpleCopier{}
}

type Splicer struct {
	tcpFD      []int
	pipe       []int
	bufferSize int
	lock       sync.Mutex
}

func (s *Splicer) CopyBytes(in, out *net.TCPConn) {
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
	inFD := int(inFile.Fd())
	defer syscall.Close(inFD)

	outFile, err := out.File()
	if err != nil {
		log.Fatal(err)
	}
	outFD := int(outFile.Fd())
	defer syscall.Close(outFD)

	var wg sync.WaitGroup
	wg.Add(2)
	go s.splice(fmt.Sprintf("from backend: %v -> %v", in.RemoteAddr(), out.RemoteAddr()), inFD, outFD, &wg)
	go s.splice(fmt.Sprintf("to backend: %v -> %v", out.RemoteAddr(), in.RemoteAddr()), outFD, inFD, &wg)
	wg.Wait()
}

func (s *Splicer) splice(direction string, inFD, outFD int, wg *sync.WaitGroup) {
	// Signal to the caller that we're done
	defer func() {
		// If a reliable delivery socket has data associated with it when a close takes place, the system continues to attempt data transfer.
		if err := syscall.Shutdown(inFD, SHUT_RDWR); err != nil {
			log.Printf("Shutdown err %v", err)
		}
		if err := syscall.Shutdown(outFD, SHUT_RDWR); err != nil {
			log.Printf("Shutdown err %v", err)
		}
		if wg != nil {
			wg.Done()
		}
	}()

	pipe := make([]int, 2)
	if err := syscall.Pipe(pipe); err != nil {
		log.Fatal(err)
	}
	defer func() {
		syscall.Close(pipe[0])
		syscall.Close(pipe[1])
	}()

	var netWrittenBytes, netReadBytes int64
	log.Printf("[%v] Splicing pipe %+v, tcpfds %+v", direction, s.pipe, s.tcpFD)
	for {
		// SPLICE_F_NONBLOCK: don't block if TCP buffer is empty
		// SPLICE_F_MOVE: directly move pages into splice buffer in kernel memory
		// SPLICE_F_MORE: just makes mysql connections slow
		log.Printf("[input (%v)] Entering input", direction)
		inBytes, err := syscall.Splice(inFD, nil, pipe[1], nil, s.bufferSize, SPLICE_F_NONBLOCK)
		if err := s.checkSpliceErr(inBytes, err, fmt.Sprintf("input (%v)", direction)); err != nil {
			log.Printf("ERROR [input (%v)] error: %v", direction, err)
			return
		}
		netReadBytes += inBytes
		log.Printf("[input (%v)] %d bytes read", direction, inBytes)

		log.Printf("[input (%v)] Entering output", direction)
		outBytes, err := syscall.Splice(pipe[0], nil, outFD, nil, s.bufferSize, SPLICE_F_NONBLOCK)
		if err := s.checkSpliceErr(inBytes, err, fmt.Sprintf("output (%v)", direction)); err != nil {
			log.Printf("ERROR [output (%v)] error: %v", direction, err)
			return
		}
		log.Printf("[output (%v)] %d bytes written, out of given input %d", direction, outBytes, inBytes)
		netWrittenBytes += outBytes
	}
	log.Printf("[%v] Spliced %d bytes read %d bytes written", direction, netWrittenBytes, netReadBytes)
}

func (s *Splicer) checkSpliceErr(bytes int64, err error, desc string) error {
	if bytes == -1 {
		if err == syscall.EINTR {
			log.Printf("[%v] ignoring EINTR with -1 bytes, syscall.%v: %v", desc, syscall.EINTR, err)
			return nil
		} else if err == syscall.EAGAIN {
			return fmt.Errorf("[%v] ignoring EAGAIN with -1 bytes, syscall.%v: %v", desc, syscall.EINTR, err)
		} else if strings.HasSuffix(err.Error(), "bad file descriptor") {
			return fmt.Errorf("[%v] bytes -1, fd closed type %T", desc, err)
		}
		return fmt.Errorf("[%v] bytes was -1 but received an unrecognized error: %v", desc, err)
	}
	if bytes == 0 {
		log.Printf("[%v] Zero bytes returned from splice", desc)
	}
	if err == nil {
		return err
	}
	switch err {
	case syscall.EINTR:
		return fmt.Errorf("[%v] Received EINTR syscall.%v: %v", desc, syscall.EINTR, err)
	case syscall.EAGAIN:
		return fmt.Errorf("[%v] Received EAGAIN syscall.%v: %v", desc, syscall.EAGAIN, err)
	default:
		return fmt.Errorf("[%v] Fail splice: %v", desc, err)
	}
}

type SimpleCopier struct{}

func (s *SimpleCopier) CopyBytes(in, out *net.TCPConn) {
	var wg sync.WaitGroup
	wg.Add(2)
	go s.copy(fmt.Sprintf("from backend: %v -> %v", in.RemoteAddr(), out.RemoteAddr()), in, out, &wg)
	go s.copy(fmt.Sprintf("to backend: %v -> %v", out.RemoteAddr(), in.RemoteAddr()), out, in, &wg)
	wg.Wait()
}

func (s *SimpleCopier) isClosedError(err error) bool {
	// A brief discussion about handling closed error here:
	// https://code.google.com/p/go/issues/detail?id=4373#c14
	// TODO: maybe create a stoppable TCP listener that returns a StoppedError
	return strings.HasSuffix(err.Error(), "use of closed network connection")
}

func (s *SimpleCopier) copy(direction string, in, out *net.TCPConn, wg *sync.WaitGroup) {
	defer func() {
		in.Close(direction)
		out.Close(direction)
		if wg != nil {
			wg.Done()
		}
	}()
	n, err := io.CopyDesc(direction, in, out)
	if err != nil {
		if s.isClosedError(err) {
			log.Printf("ERROR [%v] closed connection: %v", err)
		} else {
			log.Printf("ERROR [%v] unexpected error %v", direction, err)
		}
	}
	log.Printf("[%v] Copied %d bytes", direction, n)

}
