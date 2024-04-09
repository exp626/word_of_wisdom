package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"
	"word_of_wisdom/pkg"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	go readData(conn)

	<-ctx.Done()
	log.Println("Goodbye!")
}

func readData(conn net.Conn) {
	buffer := make([]byte, 1024)

	pow := pkg.NewEquihashPoW(48, 3, 2)

	for {
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
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Wisdom received:", string(buffer[:n]))
	}
}
