package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func echo(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Echo thread started")

	b := make([]byte, 6144)
	for {
		size, err := conn.Read(b[0:])
		if err == io.EOF {
			log.Println("Client disconnected")
		}
		if err != nil {
			log.Println("Unexpected error")
			break
		}
		now := time.Now()
		fmt.Println(now)

		if _, err := conn.Write(b[0:size]); err != nil {
			log.Fatalln("Unable to write data")
		}
	}

	log.Println("Exiting listener thread")
}

func main() {
	listener, err := net.Listen("tcp", ":20080")
	if err != nil {
		log.Fatalln("unable to bind to port")
	}
	log.Println("Listening on port 20080")

	for {
		conn, err := listener.Accept()
		log.Println("Received connection")
		if err != nil {
			log.Fatalln("Unable to accept connection")
		}
		go echo(conn)
	}
}
