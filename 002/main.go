package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sort"
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
		go handleMeans(conn)
	}
}

func handleMeans(conn net.Conn) {
	// close conn later
	defer conn.Close()
	// io.Copy(conn, conn)
	//read
	dataholder := make(map[int32]int32)
	reader := bufio.NewReader(conn)
	var timestamps []int32
	buf := make([]byte, 1024)
	messages := make(chan byte, 3000)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error: ", err)
				goto end
			}
		}
		for _, v := range buf[:n] {
			messages <- v
		}
		for len(messages) >= 9 {
			currentMsg := make([]byte, 9)
			for i := 0; i < 9; i++ {
				currentMsg[i] = <-messages
			}

			msgType := currentMsg[0]
			var data1 int32 = int32(binary.BigEndian.Uint32(currentMsg[1:5]))
			var data2 int32 = int32(binary.BigEndian.Uint32(currentMsg[5:9]))
			if msgType == 'I' { // Insert
				dataholder[data1] = data2
				timestamps = append(timestamps, data1)
			} else if msgType == 'Q' { // Query
				sort.Slice(timestamps, func(i, j int) bool {
					return timestamps[i] < timestamps[j]
				})
				var nrOfItems int64 = 0
				var sum int64 = 0
				for _, t := range timestamps {
					if t > data2 {
						break
					}
					if t < data1 {
						continue
					}
					nrOfItems++
					sum += int64(dataholder[t])
				}
				retVal := make([]byte, 4)
				var mean int64
				if nrOfItems == 0 {
					mean = 0
				} else {
					mean = sum / nrOfItems
				}
				binary.BigEndian.PutUint32(retVal, uint32(int32(mean)))
				conn.Write(retVal)

			} else {
				fmt.Println("Wrong msg type")
				goto end
			}
		}

	}

end:
	fmt.Println("Connection closed")

}
