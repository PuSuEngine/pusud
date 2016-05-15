package core

type clients []*client
type channels map[string]clients

var subscriptions = channels{}

type publishOrder struct {
	Channel string
	Data    []byte
}

var publishCn = make(chan publishOrder)

func subscribe(channel string, client *client) {
	_, ok := subscriptions[channel]

	if !ok {
		subscriptions[channel] = clients{}
	}

	subscriptions[channel] = append(subscriptions[channel], client)
}

func unsubscribe(channel string, client *client) {
	filtered := clients{}

	for _, c := range subscriptions[channel] {
		if c.UUID != client.UUID {
			filtered = append(filtered, c)
		}
	}

	subscriptions[channel] = filtered
}

func publish(channel string, data []byte) {
	_, ok := subscriptions[channel]

	if !ok {
		// Nobody listening to this channel, ignore
		return
	}

	for _, c := range subscriptions[channel] {
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
