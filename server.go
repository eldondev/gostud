package main

import (
        "crypto/rand"
        "crypto/tls"
        "log"
        "net"
)

func main() {
        cert, err := tls.LoadX509KeyPair("ca-bck.crt", "server.key")
        if err != nil {
                log.Printf("server: conn: read: %s", err)
        }
        config := tls.Config{
                ClientAuth : 0,
                Certificates: []tls.Certificate{cert},
        }
        config.Rand = rand.Reader
        service := "0.0.0.0:443"
        listener, err := tls.Listen("tcp", service, &config)
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
        buf := make([]byte, 512)
        for {
                log.Print("server: conn: waiting")
                n, err := conn.Read(buf)
                if err != nil {
                        if err != nil {
                                log.Printf("server: conn: read: %s", err)
                        }
                        break
                }


                log.Printf("server: conn: echo %q\n", string(buf[:n]))
                n, err = conn.Write(buf[:n])

                n, err = conn.Write(buf[:n])
                log.Printf("server: conn: wrote %d bytes", n)

                if err != nil {
                        log.Printf("server: write: %s", err)
                        break
                }
        }
        log.Println("server: conn: closed")
}

