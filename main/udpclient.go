package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	proto := "udp"
	name := "localhost"
	port := "5000"

	nameport := name + ":" + port

	conn, err := net.Dial(proto, nameport)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Connected: %T, %v\n", conn, conn)

	fmt.Printf("Local address: %v\n", conn.LocalAddr())
	fmt.Printf("Remote address: %v\n", conn.RemoteAddr())

	b := []byte("some string")

	cc, wrerr := conn.Write(b)

	if wrerr != nil {
		fmt.Printf("conn.Write() error: %s\n", wrerr)
	} else {
		fmt.Printf("Wrote %d bytes to socket\n", cc)
		c := make([]byte, cc+10)
		cc, rderr := conn.Read(c)
		if rderr != nil {
			fmt.Printf("conn.Read() error: %s\n", rderr)
		} else {
			fmt.Printf("Read %d bytes from socket\n", cc)
			fmt.Printf("Bytes: %q\n", string(c[0:cc]))
		}
	}

	if err = conn.Close(); err != nil {
		log.Fatal(err)
	}
}
