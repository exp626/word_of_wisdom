package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"word_of_wisdom/pkg"
)

func main() {

	addr := os.Getenv("SERVER_ADDRESS")

	fmt.Println("server address:", addr)

	if addr == "" {
		addr = "localhost:8080"
	}

	// Connect to the server
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	readData(conn)

}

func readData(conn net.Conn) {
	buffer := make([]byte, 1024)

	pow := pkg.NewEquihashPoW(48, 3, 2)

	// Read data from the client
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	challenge := buffer[:n]

	// Process and use the data (here, we'll just print it)
	fmt.Printf("Have challenge: %b\n", challenge)

	nonce, soln, err := pow.PoW(challenge)
	if err != nil {
		fmt.Println("pow error:", err)
		return
	}

	result := make([]byte, 8+len(soln)*8)

	binary.LittleEndian.PutUint64(result[:8], uint64(nonce))

	i := 8
	for _, s := range soln {
		binary.LittleEndian.PutUint64(result[i:i+8], uint64(s))
		i += 8
	}

	binary.LittleEndian.PutUint64(result, uint64(nonce))

	fmt.Printf("Have result: %b\n", result)

	_, err = conn.Write(result)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	n, err = conn.Read(buffer)
	if err != nil && !errors.Is(err, io.EOF) {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Wisdom received:", string(buffer[:n]))

}
