package models

type MessageRequest struct {
	MessageSize   uint32
	APIKey        uint16
	APIVersion    uint16
	CorrelationID uint32
	ClientID      string
	Payload       []byte
}

type MessageResponse struct {
}
