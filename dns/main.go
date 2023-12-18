package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strings"
)

var nameToIP = map[string][4]uint8{
	"habr.com":    {123, 123, 123, 123},
	"habr.ru":     {8, 8, 8, 8},
	"ru.habr.com": {1, 1, 1, 1},
}

const Address = "127.0.0.1:2053"

type Type uint16

const (
	_ Type = iota
	A
)

type Class uint16

const (
	_ Class = iota
	IN
)

type Header struct {
	PacketID uint16
	QR       uint16
	OPCODE   uint16
	AA       uint16
	TC       uint16
	RD       uint16
	RA       uint16
	Z        uint16
	RCode    uint16
	QDCount  uint16
	ANCount  uint16
	NSCount  uint16
	ARCount  uint16
}

func ReadHeader(buf []byte) Header {
	h := Header{
		PacketID: uint16(buf[0])<<8 | uint16(buf[1]),
		QR:       1,
		OPCODE:   uint16((buf[2] << 1) >> 4),
		AA:       uint16((buf[2] << 5) >> 7),
		TC:       uint16((buf[2] << 6) >> 7),
		RD:       uint16((buf[2] << 7) >> 7),
		RA:       uint16(buf[3] >> 7),
		Z:        uint16((buf[3] << 1) >> 5),
		QDCount:  uint16(buf[4])<<8 | uint16(buf[5]),
		ANCount:  uint16(buf[6])<<8 | uint16(buf[7]),
		NSCount:  uint16(buf[8])<<8 | uint16(buf[9]),
		ARCount:  uint16(buf[10])<<8 | uint16(buf[11]),
	}

	if h.OPCODE == 0 {
		h.RCode = 0
	} else {
		h.RCode = 4
	}

	return h
}

func (h Header) Encode() []byte {
	dnsHeader := make([]byte, 12)

	var flags uint16 = 0
	flags = h.QR<<15 | h.OPCODE<<11 | h.AA<<10 | h.TC<<9 | h.RD<<8 | h.RA<<7 | h.Z<<4 | h.RCode

	binary.BigEndian.PutUint16(dnsHeader[0:2], h.PacketID)
	binary.BigEndian.PutUint16(dnsHeader[2:4], flags)
	binary.BigEndian.PutUint16(dnsHeader[4:6], h.QDCount)
	binary.BigEndian.PutUint16(dnsHeader[6:8], h.ANCount)
	binary.BigEndian.PutUint16(dnsHeader[8:10], h.NSCount)
	binary.BigEndian.PutUint16(dnsHeader[10:12], h.ARCount)

	return dnsHeader
}

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", Address)
	if err != nil {
		log.Fatal("failed to resolve udp address", err)
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal("failed to to bind to address", err)
	}
	defer udpConn.Close()

	log.Printf("started server on %s", Address)

	buf := make([]byte, 512)
	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Println("failed to receive data", err)
			break
		}

		data := string(buf[:size])
		log.Printf("received %d bytes from %s: %s", size, source.String(), data)

		header := ReadHeader(buf[:12])
		log.Printf("ID: %d; QR: %d; QDCount: %d\n", header.PacketID, header.QR, header.QDCount)

		question := ReadQuestion(buf[12:])

		answer := Answer{
			Name:   question.QName,
			Type:   A,
			Class:  IN,
			TTL:    0,
			Length: net.IPv4len,
			Data:   nameToIP[question.QName],
		}

		var res bytes.Buffer
		res.Write(header.Encode())
		res.Write(question.Encode())
		res.Write(answer.Encode())

		_, err = udpConn.WriteToUDP(res.Bytes(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

func ReadQuestion(buf []byte) Question {
	start := 0
	var nameParts []string

	for len := buf[start]; len != 0; len = buf[start] {
		start++
		nameParts = append(nameParts, string(buf[start:start+int(len)]))
		start += int(len)
	}
	questionName := strings.Join(nameParts, ".")
	start++

	questionType := binary.BigEndian.Uint16(buf[start : start+2])
	questionClass := binary.BigEndian.Uint16(buf[start+2 : start+4])

	q := Question{
		QName:  questionName,
		QType:  Type(questionType),
		QClass: Class(questionClass),
	}

	return q
}

func (q Question) Encode() []byte {
	domain := q.QName
	parts := strings.Split(domain, ".")

	var buf bytes.Buffer

	for _, label := range parts {
		if len(label) > 0 {
			buf.WriteByte(byte(len(label)))
			buf.WriteString(label)
		}
	}
	buf.WriteByte(0x00)
	buf.Write(intToBytes(uint16(q.QType)))
	buf.Write(intToBytes(uint16(q.QClass)))

	return buf.Bytes()
}

type Question struct {
	QName  string
	QType  Type
	QClass Class
}

func intToBytes(n uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, n)
	return b
}

type Answer struct {
	Name   string
	Type   Type
	Class  Class
	TTL    uint32
	Length uint16
	Data   [4]uint8
}

func (a Answer) Encode() []byte {
	var rrBytes []byte

	domain := a.Name
	parts := strings.Split(domain, ".")

	for _, label := range parts {
		if len(label) > 0 {
			rrBytes = append(rrBytes, byte(len(label)))
			rrBytes = append(rrBytes, []byte(label)...)
		}
	}
	rrBytes = append(rrBytes, 0x00)

	rrBytes = append(rrBytes, intToBytes(uint16(a.Type))...)
	rrBytes = append(rrBytes, intToBytes(uint16(a.Class))...)

	time := make([]byte, 4)
	binary.BigEndian.PutUint32(time, a.TTL)

	rrBytes = append(rrBytes, time...)
	rrBytes = append(rrBytes, intToBytes(a.Length)...)

	ipBytes, err := net.IPv4(a.Data[0], a.Data[1], a.Data[2], a.Data[3]).MarshalText()
	if err != nil {
		return nil
	}

	rrBytes = append(rrBytes, ipBytes...)

	return rrBytes
}
