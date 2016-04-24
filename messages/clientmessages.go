package messages

import (
	"log"
	"encoding/json"
)


type ClientMessage struct {
	Type string `json:"type"`
}


type Hello struct {
	Type string `json:"type"`
	// TODO: Version, etc.
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


type Publish struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Content string `json:"content"`
}

func (p *Publish) ToJson() []byte {
	result, err := json.Marshal(&p)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return result
}

func NewPublish() Message {
	return &Publish{}
}


type Subscribe struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
}

func (s *Subscribe) ToJson() []byte {
	result, err := json.Marshal(&s)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return result
}

func NewSubscribe() Message {
	return &Subscribe{}
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
	return &Authorize{}
}


func init() {
	RegisterIncomingMessageType(TYPE_HELLO, NewHello)
	RegisterIncomingMessageType(TYPE_AUTHORIZE, NewAuthorize)
	RegisterIncomingMessageType(TYPE_PUBLISH, NewPublish)
	RegisterIncomingMessageType(TYPE_SUBSCRIBE, NewSubscribe)
}
