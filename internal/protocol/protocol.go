package protocol

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"gofkgo/internal/models"
	"gofkgo/internal/ports"
	"io"
	"net"
)

type ProtocolSvc struct {
}

func NewProtocolPort() ports.ProtocolPort {
	return &ProtocolSvc{}
}

func (p *ProtocolSvc) ProtocolHandler(data []byte, c net.Conn) error {
	hdr, payload, err := p.ParseMessage(data)
	if err != nil {
		return fmt.Errorf("failed to parse message: %w", err)
	}

	fmt.Printf("Parsed Header: %+v\n", *hdr)
	fmt.Printf("Payload (%d bytes): %s\n", len(payload), payload)

	// Prepare response
	respPayload := []byte(fmt.Sprintf("ACK: %s", payload))
	respSize := 4 + len(respPayload) // correlation_id (4 bytes) + payload

	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.BigEndian, uint32(respSize))  // total response size
	_ = binary.Write(&buf, binary.BigEndian, hdr.CorrelationID) // match request ID
	_, _ = buf.Write(respPayload)                               // response body

	_, err = c.Write(buf.Bytes())
	return err
}

func (p *ProtocolSvc) ReadFullMessage(reader *bufio.Reader) ([]byte, error) {
	sizeHeader, err := reader.Peek(4)
	if err != nil {
		return nil, err
	}
	messageSize := binary.BigEndian.Uint32(sizeHeader)
	totalSize := 4 + int(messageSize)

	fullMessage := make([]byte, totalSize)
	_, err = io.ReadFull(reader, fullMessage)
	return fullMessage, err
}

func (p *ProtocolSvc) ParseMessage(data []byte) (*models.MessageRequest, []byte, error) {
	buf := bytes.NewReader(data)
	var req models.MessageRequest

	if err := binary.Read(buf, binary.BigEndian, &req.MessageSize); err != nil {
		return nil, nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &req.APIKey); err != nil {
		return nil, nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &req.APIVersion); err != nil {
		return nil, nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &req.CorrelationID); err != nil {
		return nil, nil, err
	}
	var clientIDLen int16
	if err := binary.Read(buf, binary.BigEndian, &clientIDLen); err != nil {
		return nil, nil, err
	}
	if clientIDLen < 0 {
		req.ClientID = ""
	} else {
		clientIDBytes := make([]byte, clientIDLen)
		if _, err := buf.Read(clientIDBytes); err != nil {
			return nil, nil, err
		}
		req.ClientID = string(clientIDBytes)
	}

	payload := make([]byte, req.MessageSize)
	if _, err := buf.Read(payload); err != nil {
		return nil, nil, err
	}
	req.Payload = payload

	return &req, payload, nil
}
