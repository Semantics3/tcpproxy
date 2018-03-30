package main

import (
	"flag"
	"io"
	"log"
	"net"
)

func main() {
	var (
		listenAddr = flag.String("l", "", "local address to listen on")
		remoteAddr = flag.String("r", "", "remote address to dial")
	)

	flag.Parse()

	if *listenAddr == "" {
		log.Fatalf("must supply local address to listen on, -l option")
	}

	if *remoteAddr == "" {
		log.Fatalf("must supply remote address to dial, -r option")
	}

	log.Printf("Starting TCP proxy at: %v\n", *listenAddr)
	ln, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalf("listening: %v", err)
	}

	listen(ln, *remoteAddr)
}

func listen(ln net.Listener, remoteAddr string) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		go handle(conn, remoteAddr)
	}
}

func handle(conn net.Conn, remoteAddr string) {
	defer conn.Close()

	log.Printf("connected: %s", conn.RemoteAddr())
	rconn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Printf("dialing remote: %v", err)
		return
	}

	defer rconn.Close()
	errc := make(chan error)
	go copyFrom(conn, rconn, errc)
	go copyFrom(rconn, conn, errc)
	err = <-errc
	log.Printf("disconnected: %s, %v", conn.RemoteAddr(), err)
}

func copyFrom(a, b io.ReadWriter, errc chan<- error) {
	_, err := io.Copy(a, b)
	errc <- err
}
