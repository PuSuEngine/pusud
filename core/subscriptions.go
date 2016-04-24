package core

import (
	"github.com/lietu/pusud/messages"
)

type Clients []*Client
type Channels map[string]Clients

var channels = Channels{}

func Subscribe(channel string, client *Client) {
	_, ok := channels[channel]

	if !ok {
		channels[channel] = Clients{}
	}

	channels[channel] = append(channels[channel], client)
}

func Unsubscribe(channel string, client *Client) {
	filtered := Clients{}

	// Probably can't handle unsubscribing from a non-existent channel, but
	// that's ok, as this should never get called for one.
	for _, c := range (channels[channel]) {
		if c.UUID != client.UUID {
			filtered = append(filtered, c)
		}
	}

	channels[channel] = filtered
}

func Publish(message *messages.Publish) {
	_, ok := channels[message.Channel]

	if !ok {
		// Nobody listening to this channel, ignore
		return
	}

	for _, c := range(channels[message.Channel]) {
		c.SendMessage(message)
	}
}
