package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
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
		go handlePrimeTime(conn)
	}
}

func handlePrimeTime(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		var req request1
		readBytes, err := reader.ReadBytes(byte('\n'))
		fmt.Println("read bytes: ", string(readBytes))
		//readString := string(readBytes)
		//if strings.Contains(readString, "\"number\":\"") {
		//	fmt.Println("number \" problem")
		//	goto fail
		//}
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading bytes", err)
			}
			break
		}

		err = json.Unmarshal(readBytes, &req)
		var numberRaw []byte
		numberRaw = req.Number
		if bytes.ContainsAny(numberRaw, "\"") {
			fmt.Println("found a stray \"")
			goto fail
		}
		if err != nil {
			fmt.Println("error unmarshalling", err)
			goto fail
		}
		if req.Method == nil {
			fmt.Println("method == nil")
			goto fail
		}

		if *req.Method != "isPrime" {
			fmt.Println("method != isPrime")
			goto fail
		}

		//if req.Number == "" {
		//	fmt.Println("number = nil")
		//	goto fail
		//}

		var resp response1

		//check if it's a float
		numberFloat, _ := new(big.Float).SetString(string(req.Number))
		if numberFloat == nil {
			goto fail
		}
		if !numberFloat.IsInt() {
			resp.Prime = false
		} else {
			tmp, _ := numberFloat.Int(nil)
			resp.Prime = tmp.ProbablyPrime(0)
		}

		resp.Method = *req.Method

		marshal, err := json.Marshal(resp)
		if err != nil {
			fmt.Println("error while marshalling response")
		}

		marshal = append(marshal, byte('\n'))
		fmt.Println("Response is: ", string(marshal))
		conn.Write(marshal)
	}

	return
fail:
	conn.Write([]byte("{}\n"))
	return
}
