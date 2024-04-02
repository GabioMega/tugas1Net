package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:3333")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleServerConn(clientConn)
	}
}

func handleServerConn(client net.Conn) {
	//nerima pesan dari client
	defer client.Close()

	var size uint32
	err := binary.Read(client, binary.LittleEndian, &size)
	if err != nil {
		panic(err)
	}

	err = client.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		panic(err)
	}

	bytMsg := make([]byte, size)
	_, err = client.Read(bytMsg)
	netErr, ok := err.(net.Error)
	if err != nil {
		if ok && netErr.Timeout() {
			fmt.Println("Read timeout")
			return
		}
		panic(err)
	}

	msg := string(bytMsg)
	fmt.Printf("Accepted: %s\n", msg)

	var reply string
	reply = "Message has been received"

	//kirim ke client
	err = binary.Write(client, binary.LittleEndian, uint32(len(reply)))
	if err != nil {
		panic(err)
	}

	_, err = client.Write([]byte(reply))
	if err != nil {
		panic(err)
	}
}
