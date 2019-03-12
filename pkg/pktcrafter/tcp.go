package pktcrafter

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"time"
)

// Crafter is the main entity that crafts packets
type Crafter struct{}

// NewCrafter returns a new crafter
func NewCrafter() *Crafter {
	rand.Seed(time.Now().Unix())
	return &Crafter{}
}

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
	// FIN flag
	FIN = 1
	// SYN flag
	SYN = 2
	// RST flag
	RST = 4
	// PSH flag
	PSH = 8
	// ACK flag
	ACK = 16
	// URG flag
	URG = 32
)

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

// CraftSyn returns an array of bytes representing a simple SYN packet
func (c *Crafter) CraftSyn(port int) *TCPPacket {
	srcPort := random(32768, 61000)
	seq := random(0, 4294967295)

	return &TCPPacket{
		SourcePort:      uint16(srcPort),
		DestinationPort: uint16(port),
		SeqNumber:       uint32(seq),
		AckNumber:       0,
		HeaderLen:       5,
		Flags:           SYN,
		WindowSize:      65535,
	}
}

// Bytes returns the byte representation of the packet structure
func (p *TCPPacket) Bytes() ([]byte, error) {
	buf, err := p.pack()
	return buf, err
}
