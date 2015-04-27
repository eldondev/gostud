package main

import (
	"crypto/rand"
	"crypto/tls"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	cert, err := tls.LoadX509KeyPair("ca-bck.crt", "server.key")
	if err != nil {
		log.Printf("server: conn: read: %s", err)
	}
	config := tls.Config{
		ClientAuth:   0,
		Certificates: []tls.Certificate{cert},
	}
	config.Rand = rand.Reader
	listener, err := tls.Listen("tcp", os.Args[1], &config)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	log.Print("server: listening")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		defer conn.Close()
		log.Printf("server: accepted from %s", conn.RemoteAddr())
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	unwrapped, err := net.Dial("tcp", os.Args[2])
	if err != nil {
		log.Printf("server: accept: %s", err)
		return
	}
	defer unwrapped.Close()
	go io.Copy(conn, unwrapped)
	io.Copy(unwrapped, conn)
}
