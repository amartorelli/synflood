package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

// TCPPacket represents the body of a TCP packet
type TCPPacket struct {
	SourcePort      uint16
	DestinationPort uint16
	SeqNumber       uint32
	AckNumber       uint32
	HeaderLen       uint8
	Flags           uint8
	WindowSize      uint16
	Checksum        uint16
	UrgentPointer   uint16
	Options         uint32
	Data            []byte
}

// pack composes and returns the bytes representing the packet
func (p *TCPPacket) pack() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Source port
	err := binary.Write(buf, binary.BigEndian, p.SourcePort)
	if err != nil {
		return nil, err
	}
	// Destination port
	err = binary.Write(buf, binary.BigEndian, p.DestinationPort)
	if err != nil {
		return nil, err
	}
	// SeqNumber
	err = binary.Write(buf, binary.BigEndian, p.SeqNumber)
	if err != nil {
		return nil, err
	}
	// AckNumber
	err = binary.Write(buf, binary.BigEndian, p.AckNumber)
	if err != nil {
		return nil, err
	}
	// Header length + flags + Window
	hf := uint32(p.WindowSize) + uint32(p.Flags)<<16 + uint32(p.HeaderLen)<<28
	err = binary.Write(buf, binary.BigEndian, hf)
	if err != nil {
		return nil, err
	}
	// Checksum
	err = binary.Write(buf, binary.BigEndian, p.Checksum)
	if err != nil {
		return nil, err
	}
	// UrgentPointer
	err = binary.Write(buf, binary.BigEndian, p.UrgentPointer)
	if err != nil {
		return nil, err
	}
	// Options
	err = binary.Write(buf, binary.BigEndian, p.Options)
	if err != nil {
		return nil, err
	}
	// Data
	err = binary.Write(buf, binary.BigEndian, p.Data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

const (
	fin = 1
	syn = 2
	rst = 4
	psh = 8
	ack = 16
	urg = 32
)

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	interval := flag.Int("interval", 3000, "interval in milliseconds to send packets")
	host := flag.String("host", "", "the host you want to send the packets to")
	port := flag.Int("port", 0, "the port to send the packets to")
	flag.Parse()

	if *host == "" || *port == 0 {
		logrus.Fatal("host and/or port undefined")
	}

	dst, err := net.ResolveIPAddr("ip", *host)
	if err != nil {
		logrus.Fatal(err)
	}
	s, err := net.DialIP("ip4:tcp", nil, dst)
	if err != nil {
		logrus.Fatal(err)
	}

	for {
		select {
		case <-time.After(time.Duration(*interval) * time.Millisecond):
			rand.Seed(time.Now().Unix())
			srcPort := random(32768, 61000)

			packet := &TCPPacket{
				SourcePort:      uint16(srcPort),
				DestinationPort: uint16(*port),
				SeqNumber:       84710,
				AckNumber:       0,
				HeaderLen:       5,
				Flags:           syn,
				WindowSize:      65535,
			}

			buf, err := packet.pack()
			if err != nil {
				logrus.Fatal(err)
			}
			s.Write(buf)
			logrus.Printf("sent packet %v\n", packet)
		case s := <-sigs:
			logrus.Infof("received signal %s, exiting", s.String())
		}
	}

}
