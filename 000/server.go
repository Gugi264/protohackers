package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const (
	HOST = "0.0.0.0"
	PORT = "4623"
	TYPE = "tcp"
)

func main() {
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// close listener
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		fmt.Println("Got connection from: ", conn.RemoteAddr())
		go handleSmokeTest(conn)
	}
}

func handleSmokeTest(conn net.Conn) {
	// close conn later
	defer conn.Close()
	// io.Copy(conn, conn)
	//read
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error: ", err)
			}
			fmt.Printf("Writing last n bytes: %d\n", n)
			conn.Write(buf[:n])
			break
		}
		fmt.Printf("Writing bytes: %d\n", n)
		conn.Write(buf[:n])
	}
	fmt.Println("Connection closed")

}
