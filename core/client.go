package core

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/lietu/pusud/auth"
	"github.com/lietu/pusud/messages"
)

type Client struct {
	Connection  *websocket.Conn
	Request     *http.Request
	Permissions auth.Permissions
	Connected   bool
}

func (c *Client) Close() {

	if c.Connected {
		c.Connected = false
		log.Printf("Closing connection from %s", c.GetRemoteAddr())
		// TODO: De-register listeners
		c.Connection.Close()
	}
}

func (c *Client) GetRemoteAddr() string {
	return c.Request.RemoteAddr
}

func (c *Client) SendHello() {
	c.SendMessage(messages.NewHello())
}

func (c *Client) SendMessage(message messages.Message) {
	data := message.ToJson()
	log.Printf("Sending message %s to %s", data, c.GetRemoteAddr())
	c.Connection.WriteMessage(websocket.TextMessage, data)
}

func (c *Client) Authorize(message *messages.Authorize) {
	log.Printf("Client from %s authorizing with %s", c.GetRemoteAddr(), message.Authorization)

	a := GetAuthenticator()
	perms := a.GetPermissions(message.Authorization)

	if len(perms) == 0 {
		// No permissions granted -> invalid authorization
		c.SendMessage(messages.NewGenericMessage(messages.TYPE_AUTHORIZATION_FAILED))
		c.Close()
		return
	}

	for channel, perm := range (perms) {
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
}

func (c *Client) ReadMessage(content []byte) {
	m := messages.NewMessageFromContent(content)

	if a, ok := m.(*messages.Authorize); ok {
		c.Authorize(a)
	} else {
		// Unknown message type
		c.SendMessage(messages.NewGenericMessage(messages.TYPE_UNKNOWN_MESSAGE_RECEIVED))
		c.Close()
	}
}

func (c *Client) Handle() {
	for {
		_, message, err := c.Connection.ReadMessage()

		log.Printf("%s", message)

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

func NewClient(conn *websocket.Conn, req *http.Request) *Client {
	c := Client{}
	c.Connection = conn
	c.Request = req
	c.Permissions = auth.Permissions{}
	c.Connected = true

	return &c
}
