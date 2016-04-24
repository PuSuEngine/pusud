package messages

import (
	"log"
	"encoding/json"
)

const TYPE_HELLO = "hello"
const TYPE_AUTHORIZE = "authorize"
const TYPE_UNKNOWN_MESSAGE_RECEIVED = "unknown_message_received"
const TYPE_AUTHORIZATION_FAILED = "authorization_failed"


type ClientMessage struct {
	Type string `json:"type"`
}


type Hello struct {
	Type string `json:"type"`
}

func (h *Hello) ToJson() []byte {
	result, err := json.Marshal(&h)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return result
}

func NewHello() Message {
	h := Hello{}
	h.Type = TYPE_HELLO

	return &h
}


type Authorize struct {
	Type          string `json:"type"`
	Authorization string `json:"authorization"`
}

func (m *Authorize) ToJson() []byte {
	result, err := json.Marshal(&m)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return result
}

func NewAuthorize() Message {
	a := Authorize{}
	a.Type = TYPE_AUTHORIZE

	return &a
}


func init() {
	RegisterIncomingMessageType(TYPE_HELLO, NewHello)
	RegisterIncomingMessageType(TYPE_AUTHORIZE, NewAuthorize)
}
