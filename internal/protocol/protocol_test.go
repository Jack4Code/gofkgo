package protocol

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"testing"
)

// frame builds a length-prefixed request with a v1 header.
func frame(apiKey, apiVersion int16, corrID int32, clientID *string, body []byte) []byte {
	var msg bytes.Buffer
	binary.Write(&msg, binary.BigEndian, apiKey)
	binary.Write(&msg, binary.BigEndian, apiVersion)
	binary.Write(&msg, binary.BigEndian, corrID)
	if clientID == nil {
		binary.Write(&msg, binary.BigEndian, int16(-1))
	} else {
		binary.Write(&msg, binary.BigEndian, int16(len(*clientID)))
		msg.WriteString(*clientID)
	}
	msg.Write(body)

	var out bytes.Buffer
	binary.Write(&out, binary.BigEndian, int32(msg.Len()))
	out.Write(msg.Bytes())
	return out.Bytes()
}

func ptr(s string) *string { return &s }

func TestParseMessage(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		apiKey   int16
		version  int16
		corrID   int32
		clientID string
		payload  []byte
	}{
		{"empty body", frame(18, 0, 1, ptr("kcat"), nil), 18, 0, 1, "kcat", []byte{}},
		{"with body", frame(0, 3, 7, ptr("kcat"), []byte("hello")), 0, 3, 7, "kcat", []byte("hello")},
		{"null client id", frame(3, 1, 42, nil, []byte{1, 2}), 3, 1, 42, "", []byte{1, 2}},
		{"empty client id", frame(3, 1, 42, ptr(""), nil), 3, 1, 42, "", []byte{}},
		{"negative correlation id", frame(18, 0, -5, ptr("c"), nil), 18, 0, -5, "c", []byte{}},
	}

	p := &ProtocolSvc{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, payload, err := p.ParseMessage(tt.data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if req.APIKey != tt.apiKey || req.APIVersion != tt.version || req.CorrelationID != tt.corrID {
				t.Errorf("header = {key %d, ver %d, corr %d}, want {%d, %d, %d}",
					req.APIKey, req.APIVersion, req.CorrelationID, tt.apiKey, tt.version, tt.corrID)
			}
			if req.ClientID != tt.clientID {
				t.Errorf("client id = %q, want %q", req.ClientID, tt.clientID)
			}
			if !bytes.Equal(payload, tt.payload) {
				t.Errorf("payload = %v, want %v", payload, tt.payload)
			}
		})
	}
}

func TestParseMessageErrors(t *testing.T) {
	good := frame(18, 0, 1, ptr("kcat"), nil)

	// Client id length claims more bytes than are present.
	shortClientID := []byte{0, 0, 0, 12, 0, 18, 0, 0, 0, 0, 0, 1, 0, 9, 'a', 'b'}

	tests := []struct {
		name string
		data []byte
		want error
	}{
		{"empty", nil, io.EOF},
		{"size only", []byte{0, 0, 0, 0}, io.EOF},
		{"truncated header", []byte{0, 0, 0, 3, 0, 18, 0}, io.ErrUnexpectedEOF},
		{"short client id", shortClientID, io.ErrUnexpectedEOF},
		{"size mismatch", good[:len(good)-1], ErrSizeMismatch},
	}

	p := &ProtocolSvc{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := p.ParseMessage(tt.data)
			if !errors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestReadFullMessage(t *testing.T) {
	p := &ProtocolSvc{}

	t.Run("back to back messages", func(t *testing.T) {
		a := frame(18, 0, 1, ptr("kcat"), nil)
		b := frame(0, 0, 2, ptr("kcat"), []byte("hello"))
		r := bufio.NewReader(bytes.NewReader(append(append([]byte{}, a...), b...)))

		for i, want := range [][]byte{a, b} {
			got, err := p.ReadFullMessage(r)
			if err != nil {
				t.Fatalf("message %d: %v", i, err)
			}
			if !bytes.Equal(got, want) {
				t.Fatalf("message %d = %v, want %v", i, got, want)
			}
		}
		if _, err := p.ReadFullMessage(r); err != io.EOF {
			t.Fatalf("after last message err = %v, want EOF", err)
		}
	})

	tests := []struct {
		name string
		data []byte
		want error
	}{
		{"too large", []byte("GET / HTTP/1.1\r\n"), ErrMessageTooLarge},
		{"negative size", []byte{0xff, 0xff, 0xff, 0xff}, ErrNegativeSize},
		{"truncated body", []byte{0, 0, 0, 10, 1, 2, 3}, io.ErrUnexpectedEOF},
		{"truncated size", []byte{0, 0}, io.EOF},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := p.ReadFullMessage(bufio.NewReader(bytes.NewReader(tt.data)))
			if !errors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
			}
		})
	}
}
