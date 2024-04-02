package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

func sendtoServer(msg string) {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:3333", 3*time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	//kirim msg ke server
	err = binary.Write(conn, binary.LittleEndian, uint32(len(msg)))
	if err != nil {
		panic(err)
	}

	err = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		panic(err)
	}

	_, err = conn.Write([]byte(msg))
	netErr, oke := err.(net.Error)
	if err != nil {
		if oke && netErr.Timeout() {
			fmt.Println("Write timeout")
			return
		}
		panic(err)
	}

	//nerima reply dari server
	var size uint32
	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		panic(err)
	}
	err = binary.Read(conn, binary.LittleEndian, &size)
	netErr, ok := err.(net.Error)
	if err != nil {
		if ok && netErr.Timeout() {
			fmt.Println("Read timeout")
			return
		}
		panic(err)
	}
	err = conn.SetReadDeadline(time.Time{})
	if err != nil {
		panic(err)
	}

	bytReply := make([]byte, size)
	_, err = conn.Read(bytReply)
	if err != nil {
		panic(err)
	}
	fmt.Printf("reply: %s\n", string(bytReply))
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("1. Send a message")
		fmt.Println("2. Exit")
		fmt.Print(">> ")
		scanner.Scan()
		input := scanner.Text()
		if input == "1" {
			var msg string
			for {
				fmt.Print("Enter a message: ")
				scanner.Scan()
				msg = scanner.Text()
				if len(msg) > 0 {
					break
				}
			}
			sendtoServer(msg)
		} else if input == "2" {
			break
		}
	}
}
