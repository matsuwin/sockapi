package main

import (
	"fmt"
	"os"
	"sockapi"
	"time"
)

// var proto, bindAddr = "unix", "/tmp/echo.sock"
var proto, bindAddr = "tcp", ":3000"

// Socket C/S: socket > connect
func client() {
	sa, err := sockapi.SocketAddrDo(proto, bindAddr, sockapi.ClientMode)
	defer sa.Close()
	if err != nil {
		panic(err)
	}

	// write request
	if _, err = sa.Write([]byte("ping")); err != nil {
		panic(err)
	}

	// read response
	buf, err := sa.Read(1024)
	if err != nil {
		panic(err)
	}
	fmt.Printf("response: %s\n", buf)
}

// Socket C/S: socket > bind > listen
func main() {
	if err := os.RemoveAll(bindAddr); err != nil {
		panic(err)
	}
	sa, err := sockapi.SocketAddrDo(proto, bindAddr, sockapi.ServerMode)
	defer sa.Close()
	if err != nil {
		panic(err)
	}
	go accept(sa)
	time.Sleep(time.Second)

	t := time.NewTicker(time.Second)
	for range t.C {
		client()
	}
}

func accept(sa *sockapi.SocketAddr) {
	for {
		nsa, _, err := sa.Accept()
		if err != nil {
			panic(err)
		}

		// read request
		payload, err := nsa.Read(1024)
		if err != nil {
			panic(err)
		}
		handler(nsa, payload)
	}
}

func handler(sa *sockapi.SocketAddr, payload []byte) {
	defer sa.Close()
	fmt.Printf("request: %s\n", payload)

	// write response
	if _, err := sa.Write([]byte("pong")); err != nil {
		panic(err)
	}
}
