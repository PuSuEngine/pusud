package messages

import (
	"encoding/json"
	"log"
)

type Message interface {
	ToJson() []byte
}

type MessageConstructor func() Message

type GenericMessage struct {
	Type string `json:"type"`
}

func (um *GenericMessage) ToJson() []byte {
	result, err := json.Marshal(&um)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return result
}

func NewGenericMessage(messageType string) Message {
	g := GenericMessage{}
	g.Type = messageType
	return &g
}
