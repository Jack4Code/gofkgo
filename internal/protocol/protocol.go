package protocol

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"gofkgo/internal/models"
	"gofkgo/internal/ports"
	"io"
	"net"
)

// MaxMessageSize matches Kafka's default socket.request.max.bytes.
const MaxMessageSize = 100 * 1024 * 1024

var (
	ErrMessageTooLarge = errors.New("message size exceeds limit")
	ErrNegativeSize    = errors.New("negative message size")
	ErrSizeMismatch    = errors.New("message size does not match frame length")
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
	_ = binary.Write(&buf, binary.BigEndian, int32(respSize))   // total response size
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
	messageSize := int32(binary.BigEndian.Uint32(sizeHeader))
	if messageSize < 0 {
		return nil, ErrNegativeSize
	}
	if messageSize > MaxMessageSize {
		return nil, fmt.Errorf("%w: %d bytes", ErrMessageTooLarge, messageSize)
	}

	fullMessage := make([]byte, 4+int(messageSize))
	if _, err := io.ReadFull(reader, fullMessage); err != nil {
		return nil, err
	}
	return fullMessage, nil
}

// ParseMessage parses a request header (v1: no tagged fields) and returns
// whatever follows it as the payload.
func (p *ProtocolSvc) ParseMessage(data []byte) (*models.MessageRequest, []byte, error) {
	buf := bytes.NewReader(data)
	var req models.MessageRequest

	if err := binary.Read(buf, binary.BigEndian, &req.MessageSize); err != nil {
		return nil, nil, err
	}
	if int(req.MessageSize) != len(data)-4 {
		return nil, nil, fmt.Errorf("%w: header says %d, got %d", ErrSizeMismatch, req.MessageSize, len(data)-4)
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
	if clientIDLen >= 0 {
		clientIDBytes := make([]byte, clientIDLen)
		if _, err := io.ReadFull(buf, clientIDBytes); err != nil {
			return nil, nil, err
		}
		req.ClientID = string(clientIDBytes)
	}

	offset := len(data) - buf.Len()
	req.Payload = data[offset:]

	return &req, req.Payload, nil
}
