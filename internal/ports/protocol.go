package ports

import (
	"bufio"
	"gofkgo/internal/models"
	"net"
)

type ProtocolPort interface {
	ProtocolHandler([]byte, net.Conn) error
	ReadFullMessage(reader *bufio.Reader) ([]byte, error)
	ParseMessage(data []byte) (*models.MessageRequest, []byte, error)
}
