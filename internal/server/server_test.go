package server

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"testing"
	"time"

	"gofkgo/internal/protocol"
)

func TestHandleConnection(t *testing.T) {
	client, srv := net.Pipe()
	defer client.Close()

	s := &Server{Handler: protocol.NewProtocolPort()}
	go s.handleConnection(srv)

	client.SetDeadline(time.Now().Add(2 * time.Second))

	// Two requests on one connection, the first with an empty body
	// (like ApiVersions v0). Both should get a response.
	for _, tc := range []struct {
		corrID int32
		body   string
	}{
		{1, ""},
		{2, "hello"},
	} {
		var msg bytes.Buffer
		binary.Write(&msg, binary.BigEndian, int16(18))
		binary.Write(&msg, binary.BigEndian, int16(0))
		binary.Write(&msg, binary.BigEndian, tc.corrID)
		binary.Write(&msg, binary.BigEndian, int16(4))
		msg.WriteString("test")
		msg.WriteString(tc.body)

		req := binary.BigEndian.AppendUint32(nil, uint32(msg.Len()))
		req = append(req, msg.Bytes()...)
		if _, err := client.Write(req); err != nil {
			t.Fatalf("write: %v", err)
		}

		var size, corrID int32
		if err := binary.Read(client, binary.BigEndian, &size); err != nil {
			t.Fatalf("read size: %v", err)
		}
		if err := binary.Read(client, binary.BigEndian, &corrID); err != nil {
			t.Fatalf("read correlation id: %v", err)
		}
		body := make([]byte, size-4)
		if _, err := io.ReadFull(client, body); err != nil {
			t.Fatalf("read body: %v", err)
		}

		if corrID != tc.corrID {
			t.Errorf("correlation id = %d, want %d", corrID, tc.corrID)
		}
		if want := "ACK: " + tc.body; string(body) != want {
			t.Errorf("body = %q, want %q", body, want)
		}
	}
}
