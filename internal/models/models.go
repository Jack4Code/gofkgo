package models

type MessageRequest struct {
	MessageSize   int32
	APIKey        int16
	APIVersion    int16
	CorrelationID int32
	ClientID      string
	Payload       []byte
}

type MessageResponse struct {
}
