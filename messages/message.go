package messages

import (
	"log"
	"encoding/json"
)

var messageTypes = map[string]MessageConstructor{}

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

func NewMessageFromContent(content []byte) Message {
	um := GenericMessage{}
	json.Unmarshal(content, &um)

	var constructor MessageConstructor
	constructor, _ = messageTypes[um.Type]

	if constructor == nil {
		return &um
	}

	m := constructor()
	json.Unmarshal(content, &m)

	return m
}

func RegisterIncomingMessageType(name string, constructor MessageConstructor) {
	messageTypes[name] = constructor
}