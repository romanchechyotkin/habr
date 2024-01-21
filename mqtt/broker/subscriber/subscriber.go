package subscriber

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

type Options struct {
	Topic   string
	Address string
}

type Subscriber struct {
	conn    net.Conn
	topic   string
	msgChan chan []byte
}

func New(opts *Options) (*Subscriber, error) {
	conn, err := net.Dial("tcp", opts.Address)
	if err != nil {
		log.Fatalf("failed to dial address: %s, error: %v\n", opts.Address, err)
	}

	s := &Subscriber{
		conn:    conn,
		topic:   opts.Topic,
		msgChan: make(chan []byte),
	}

	var buf bytes.Buffer

	header := make([]byte, 2)
	binary.BigEndian.PutUint16(header, 0b01000000)

	buf.Write(header)
	buf.WriteString(s.topic)

	_, err = s.conn.Write(buf.Bytes())
	if err != nil {
		log.Printf("failed to create subscriber, error: %v\n", err)
		return nil, err
	}

	go s.readMessages()

	return s, nil
}

func (sub *Subscriber) readMessages() {
	defer sub.conn.Close()

	buf := make([]byte, 1024)
	for {
		size, err := sub.conn.Read(buf)
		if err != nil {
			log.Printf("failed to read message, error: %v\n", err)
			return
		}
		sub.msgChan <- buf[:size]
	}
}

func (sub *Subscriber) GetMsgChannel() <-chan []byte {
	return sub.msgChan
}
