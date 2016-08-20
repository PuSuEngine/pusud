package core

import (
	"sync"
)

type clients []*client
type channels map[string]clients

var subscriptionMutex = &sync.Mutex{}
var subscriptions = channels{}

type publishOrder struct {
	Channel string
	Data    []byte
}

var publishCn = make(chan publishOrder)

func subscribe(channel string, client *client) {
	subscriptionMutex.Lock()
	defer subscriptionMutex.Unlock()

	_, ok := subscriptions[channel]

	if !ok {
		subscriptions[channel] = clients{}
	}

	subscriptions[channel] = append(subscriptions[channel], client)
}

func unsubscribe(channel string, client *client) {
	subscriptionMutex.Lock()
	defer subscriptionMutex.Unlock()

	filtered := clients{}

	for _, c := range subscriptions[channel] {
		if c.UUID != client.UUID {
			filtered = append(filtered, c)
		}
	}

	subscriptions[channel] = filtered
}

func getSubscriptions(channel string) (list clients, ok bool) {
	subscriptionMutex.Lock()
	defer subscriptionMutex.Unlock()

	list, ok = subscriptions[channel]
	return
}

func publish(channel string, data []byte) {
	clients, ok := getSubscriptions(channel)

	if !ok {
		// Nobody listening to this channel, ignore
		return
	}

	for _, c := range clients {
		go c.Send(data)
	}
}

func init() {
	go func() {
		for {
			select {
			case order := <-publishCn:
				publish(order.Channel, order.Data)
			}
		}
	}()
}
