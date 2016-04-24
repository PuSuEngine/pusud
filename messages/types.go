package messages

import (
	"encoding/json"
)

const TYPE_HELLO = "hello"
const TYPE_AUTHORIZE = "authorize"
const TYPE_AUTHORIZATION_OK = "authorization_ok"
const TYPE_PUBLISH = "publish"
const TYPE_SUBSCRIBE = "subscribe"
const TYPE_SUBSCRIBE_OK = "subscribe_ok"
const TYPE_UNKNOWN_MESSAGE_RECEIVED = "unknown_message_received"
const TYPE_AUTHORIZATION_FAILED = "authorization_failed"
const TYPE_PERMISSION_DENIED = "permission_denied"


var messageTypes = map[string]MessageConstructor{}

func RegisterIncomingMessageType(name string, constructor MessageConstructor) {
	messageTypes[name] = constructor
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
