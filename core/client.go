package core

import (
	"github.com/lietu/pusud/auth"
	"github.com/lietu/pusud/messages"

	"github.com/gorilla/websocket"
	"github.com/nu7hatch/gouuid"
	"log"
	"net/http"
)

const debug = false

type permissionCache map[string]bool

type client struct {
	UUID          string
	Connection    *websocket.Conn
	Request       *http.Request
	Permissions   auth.Permissions
	Connected     bool
	Subscriptions []string
	Write         permissionCache
}

var connectedClients int64 = 0

func (c *client) Close() {

	if c.Connected {
		connectedClients--
		c.Connected = false

		if debug {
			log.Printf("Closing connection from %s", c.GetRemoteAddr())
		}

		for _, channel := range c.Subscriptions {
			unsubscribe(channel, c)
		}

		c.Connection.Close()
	}
}

func (c *client) GetRemoteAddr() string {
	return c.Request.RemoteAddr
}

func (c *client) GetPermissions(channel string) (read bool, write bool) {
	return auth.GetChannelPermissions(channel, c.Permissions)
}

func (c *client) SendHello() {
	c.SendMessage(messages.NewHello())
}

func (c *client) SendMessage(message messages.Message) {
	data := message.ToJson()
	if debug {
		log.Printf("Sending message %s to %s", data, c.GetRemoteAddr())
	}

	c.SendRaw(data)
}

func (c *client) SendRaw(data []byte) {
	c.Connection.WriteMessage(websocket.TextMessage, data)
}

func (c *client) Authorize(message *messages.Authorize) {
	if debug {
		log.Printf("Client from %s authorizing with %s", c.GetRemoteAddr(), message.Authorization)
	}

	a := getAuthenticator()
	perms := a.GetPermissions(message.Authorization)

	if len(perms) == 0 {
		// No permissions granted -> invalid authorization
		c.SendMessage(messages.NewGenericMessage(messages.TYPE_AUTHORIZATION_FAILED))
		c.Close()
		return
	}

	for channel, perm := range perms {
		old, ok := c.Permissions[channel]

		if ok {
			c.Permissions[channel] = &auth.Permission{
				old.Read || perm.Read,
				old.Write || perm.Write,
			}
		} else {
			c.Permissions[channel] = perm
		}
	}

	c.SendMessage(messages.NewGenericMessage(messages.TYPE_AUTHORIZATION_OK))
}

func (c *client) Publish(message *messages.Publish, data []byte) {
	if debug {
		log.Printf("Client from %s publishing %s to %s", c.GetRemoteAddr(), message.Content, message.Channel)
	}

	// We only need to check write permission once
	if _, ok := c.Write[message.Channel]; !ok {
		_, write := c.GetPermissions(message.Channel)
		if !write {
			c.SendMessage(messages.NewGenericMessage(messages.TYPE_PERMISSION_DENIED))
			c.Close()
			return
		}
		c.Write[message.Channel] = true
	}

	publishCn <- publishOrder{message.Channel, data}
}

func (c *client) Subscribe(message *messages.Subscribe) {
	if debug {
		log.Printf("Client from %s subscribing to %s", c.GetRemoteAddr(), message.Channel)
	}

	// Ignore double-subscription
	if c.IsSubscribed(message.Channel) {
		return
	}

	read, _ := c.GetPermissions(message.Channel)

	if !read {
		c.SendMessage(messages.NewGenericMessage(messages.TYPE_PERMISSION_DENIED))
		c.Close()
		return
	}

	c.Subscriptions = append(c.Subscriptions, message.Channel)
	subscribe(message.Channel, c)
	c.SendMessage(messages.NewGenericMessage(messages.TYPE_SUBSCRIBE_OK))
}

func (c *client) IsSubscribed(channel string) bool {
	for _, cn := range c.Subscriptions {
		if cn == channel {
			return true
		}
	}

	return false
}

func (c *client) ReadMessage(content []byte) {
	m := messages.NewMessageFromContent(content)

	if a, ok := m.(*messages.Authorize); ok {
		c.Authorize(a)
	} else if p, ok := m.(*messages.Publish); ok {
		c.Publish(p, content)
	} else if s, ok := m.(*messages.Subscribe); ok {
		c.Subscribe(s)
	} else {
		// Unknown message type
		c.SendMessage(messages.NewGenericMessage(messages.TYPE_UNKNOWN_MESSAGE_RECEIVED))
		c.Close()
	}
}

func (c *client) Handle() {
	for {
		_, message, err := c.Connection.ReadMessage()

		if err != nil {
			if err.Error() == "websocket: close 1001 " {
				log.Printf("Client from %s disconnected", c.GetRemoteAddr())
			} else {
				log.Printf("Client from %s error: %s", c.GetRemoteAddr(), err.Error())
			}
			c.Close()
			break
		} else {
			c.ReadMessage(message)
		}
	}
}

func newClient(conn *websocket.Conn, req *http.Request) *client {
	id, err := uuid.NewV4()

	if err != nil {
		log.Fatalf("UUID error: %s", err.Error())
	}

	c := client{}
	c.UUID = id.String()
	c.Connection = conn
	c.Request = req
	c.Permissions = auth.Permissions{}
	c.Connected = true
	c.Subscriptions = []string{}
	c.Write = permissionCache{}

	connectedClients++

	return &c
}
