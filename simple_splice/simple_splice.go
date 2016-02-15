package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"syscall"
)

var (
	inFile  = flag.String("in-file", "foo", "path to input file")
	outFile = flag.String("out-file", "bar", "path to input file")
)

const (
	// These constants are pulled from kernel source.
	SPLICE_F_MOVE     = 0x01
	SPLICE_F_MORE     = 0x04
	SPLICE_F_NONBLOCK = 002
)

func fileSize(f *os.File) int64 {
	info, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	return info.Size()
}

func isClosedError(err error) bool {
	// A brief discussion about handling closed error here:
	// https://code.google.com/p/go/issues/detail?id=4373#c14
	// TODO: maybe create a stoppable TCP listener that returns a StoppedError
	return strings.HasSuffix(err.Error(), "bad file descriptor")
}

func checkSpliceErr(bytes int64, err error, desc string) {
	if bytes == -1 {
		if err == syscall.EINTR {
			log.Printf("[%v] ignoring EINTR with -1 bytes, syscall.%v: %v", desc, syscall.EINTR, err)
			return
		} else if err == syscall.EAGAIN {
			log.Fatalf("[%v] ignoring EAGAIN with -1 bytes, syscall.%v: %v", desc, syscall.EINTR, err)
		} else if isClosedError(err) {
			log.Fatalf("[%v] bytes -1, fd closed type %T", desc, err)
		}
		log.Fatalf("[%v] bytes was -1 but received an unrecognized error: %v", desc, err)
	}
	if bytes == 0 {
		log.Printf("[%v] Zero bytes returned from splice", desc)
	}
	if err == nil {
		return
	}
	switch err {
	case syscall.EINTR:
		log.Fatalf("[%v] Received EINTR syscall.%v: %v", desc, syscall.EINTR, err)
	case syscall.EAGAIN:
		log.Fatalf("[%v] Received EAGAIN syscall.%v: %v", desc, syscall.EAGAIN, err)
	default:
		log.Fatalf("[%v] Fail splice: %v", desc, err)
	}
}

func main() {
	log.Printf("Input file %v", *inFile)
	in, err := os.Open(*inFile)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()
	inFD := int(in.Fd())

	log.Printf("Output file %v", *outFile)
	out, err := os.Create(*outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	outFD := int(out.Fd())

	pipe := make([]int, 2)
	if err := syscall.Pipe(pipe); err != nil {
		log.Fatal(err)
	}
	defer syscall.Close(pipe[0])
	defer syscall.Close(pipe[1])

	flags := SPLICE_F_NONBLOCK
	blockSize := 1
	inFileSize := fileSize(in)
	var netRead int64

	for {
		inBytes, err := syscall.Splice(inFD, nil, pipe[1], nil, blockSize, flags)
		checkSpliceErr(inBytes, err, "input")
		netRead += inBytes
		log.Printf("[input] %d bytes read, %d remaining", inBytes, inFileSize-netRead)

		if err := syscall.SetNonblock(inFD, true); err != nil {
			log.Fatalf("Unable to close in fd")
		}

		outBytes, err := syscall.Splice(pipe[0], nil, outFD, nil, blockSize, flags)
		checkSpliceErr(outBytes, err, "output")
		log.Printf("[output] %d bytes written, out of given input %d", outBytes, inBytes)
	}
}
