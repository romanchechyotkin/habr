package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

const Address = "127.0.0.1:1883"
const HeaderSize = 2

var subscribers = make(map[net.Conn]string)

func main() {
	log.Println("mqtt broker start")

	listener, err := net.Listen("tcp", Address)
	if err != nil {
		log.Fatalf("failed to listen address: %s, error: %v\n", Address, err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection, error: %v\n", err)
			continue
		}

		go listenConnection(conn)
	}
}

func listenConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		delete(subscribers, conn)
		log.Printf("connection %s closed\n", conn.RemoteAddr().String())
	}()

	log.Printf("incoming connection from %s\n", conn.RemoteAddr().String())

	buf := make([]byte, 1024)
	for {
		size, err := conn.Read(buf)
		if err != nil {
			log.Printf("failed to read message, error: %v\n", err)
			return
		}

		log.Printf("got message %q from %s", buf[:size], conn.RemoteAddr().String())

		number := binary.BigEndian.Uint16(buf[:HeaderSize])
		pub := number >> 7
		sub := number << 9 >> 15

		if sub == 1 {
			topic := string(buf[HeaderSize:size])
			subscribers[conn] = topic
			log.Println("got subscriber", topic)
			continue
		}

		if pub == 1 {
			idx := bytes.LastIndexByte(buf[:size], 0x00)
			pubTopic := string(buf[HeaderSize:idx])

			for sub, topic := range subscribers {
				if pubTopic == topic {
					_, err := sub.Write(buf[idx:size])
					if err != nil {
						log.Printf("failed to write message to %s, error: %v\n", conn.RemoteAddr().String(), err)
						continue
					}
				}
			}
		}
	}

}
