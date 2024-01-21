package publisher

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

type Publisher struct {
	conn  net.Conn
	topic string
}

func New(opts *Options) (*Publisher, error) {
	conn, err := net.Dial("tcp", opts.Address)
	if err != nil {
		log.Printf("failed to dial address: %s, error: %v\n", opts.Address, err)
		return nil, err
	}

	p := &Publisher{
		conn:  conn,
		topic: opts.Topic,
	}

	return p, nil
}

func (pub *Publisher) Write(msg []byte) error {
	var buf bytes.Buffer

	header := make([]byte, 2)
	binary.BigEndian.PutUint16(header, 0b10000000)

	buf.Write(header)
	buf.WriteString(pub.topic)
	buf.Write([]byte{0x00})
	buf.Write(msg)

	_, err := pub.conn.Write(buf.Bytes())
	if err != nil {
		log.Printf("failed to write message to %s, error: %v\n", pub.conn.RemoteAddr().String(), err)
		pub.conn.Close()
		return err
	}

	return nil
}
