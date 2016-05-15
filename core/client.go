package core

import (
	"github.com/PuSuEngine/pusud/auth"
	"github.com/PuSuEngine/pusud/messages"

	"github.com/gorilla/websocket"
	"github.com/nu7hatch/gouuid"
	"log"
	"sync"
	"net/http"
)

const debug = false

// Buffer up to this many messages going out
const OUTGOING_BUFFER = 100
// Buffer up to this many messages coming in
const INCOMING_BUFFER = 100

var readCounter = 0
var writeCounter = 0

type permissionCache map[string]bool

type client struct {
	UUID          string
	connection    *websocket.Conn
	request       *http.Request
	permissions   auth.Permissions
	connected     bool
	subscriptions []string
	write         permissionCache
	outgoing      chan []byte
	incoming      chan []byte
	stop          chan bool
	stopping      bool
	mutex		  *sync.Mutex
}

var connectedClients int64 = 0

func (c *client) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.connected {
		c.connected = false
		connectedClients--

		if !c.stopping {
			c.stop <- true
		}

		// Close the channels, throw away any extra messages that might be coming our way
		close(c.outgoing)
		close(c.incoming)
		close(c.stop)

		if debug {
			log.Printf("Closing connection from %s", c.GetRemoteAddr())
		}

		for _, channel := range c.subscriptions {
			unsubscribe(channel, c)
		}

		c.connection.Close()
	}
}

func (c *client) GetRemoteAddr() string {
	return c.request.RemoteAddr
}

func (c *client) GetPermissions(channel string) (read bool, write bool) {
	return auth.GetChannelPermissions(channel, c.permissions)
}

func (c *client) SendHello() {
	c.SendMessage(messages.NewHello())
}

func (c *client) SendMessage(message messages.Message) {
	data := message.ToJson()
	if debug {
		log.Printf("Sending message %s to %s", data, c.GetRemoteAddr())
	}

	c.Send(data)
}

func (c *client) sendRaw(data []byte) {
	c.connection.WriteMessage(websocket.TextMessage, data)
}

func (c *client) authorize(message *messages.Authorize) {
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
		old, ok := c.permissions[channel]

		if ok {
			c.permissions[channel] = &auth.Permission{
				old.Read || perm.Read,
				old.Write || perm.Write,
			}
		} else {
			c.permissions[channel] = perm
		}
	}

	c.SendMessage(messages.NewGenericMessage(messages.TYPE_AUTHORIZATION_OK))
}

func (c *client) publish(message *messages.Publish, data []byte) {
	if debug {
		log.Printf("Client from %s publishing %s to %s", c.GetRemoteAddr(), message.Content, message.Channel)
	}

	// We only need to check write permission once
	if _, ok := c.write[message.Channel]; !ok {
		_, write := c.GetPermissions(message.Channel)
		if !write {
			c.SendMessage(messages.NewGenericMessage(messages.TYPE_PERMISSION_DENIED))
			c.Close()
			return
		}
		c.write[message.Channel] = true
	}

	publishCn <- publishOrder{message.Channel, data}
}

func (c *client) subscribe(message *messages.Subscribe) {
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

	c.subscriptions = append(c.subscriptions, message.Channel)
	subscribe(message.Channel, c)
	c.SendMessage(messages.NewGenericMessage(messages.TYPE_SUBSCRIBE_OK))
}

func (c *client) unsubscribe(message *messages.Unsubscribe) {
	if debug {
		log.Printf("Client from %s unsubscribing from %s", c.GetRemoteAddr(), message.Channel)
	}

	if !c.IsSubscribed(message.Channel) {
		return
	}

	var filtered []string

	// Remove channel from subscriptions for that client
	for _, c := range c.subscriptions {
		if message.Channel != c {
			filtered = append(filtered, c)
		}
	}

	c.subscriptions = filtered
	unsubscribe(message.Channel, c)
}

func (c *client) IsSubscribed(channel string) bool {
	for _, cn := range c.subscriptions {
		if cn == channel {
			return true
		}
	}

	return false
}

func (c *client) readMessage(content []byte) {
	m := messages.NewMessageFromContent(content)

	if a, ok := m.(*messages.Authorize); ok {
		c.authorize(a)
	} else if p, ok := m.(*messages.Publish); ok {
		c.publish(p, content)
	} else if s, ok := m.(*messages.Subscribe); ok {
		c.subscribe(s)
	} else if s, ok := m.(*messages.Unsubscribe); ok {
		c.unsubscribe(s)
	} else {
		// Unknown message type
		c.SendMessage(messages.NewGenericMessage(messages.TYPE_UNKNOWN_MESSAGE_RECEIVED))
		c.Close()
	}
}

func (c *client) Send(data []byte) {
	select {
	case c.outgoing <- data:
		// Message sent to outgoing queue
	default:
		if c.stopping || !c.connected {
			return
		}

		log.Printf("Client from %s filled outgoing message queue, dropping connection.", c.GetRemoteAddr())
		c.stop <- true
	}

}

func (c *client) handleChannels() {
	for {
		select {
		case <- c.stop:
			if !c.stopping {
				c.stopping = true
				c.Close()
			}
			return
		case msg := <- c.outgoing:
			writeCounter++
			c.sendRaw(msg)
		case msg := <- c.incoming:
			readCounter++
			c.readMessage(msg)
		}
	}
}

func (c *client) Handle() {
	go c.handleChannels()

	for {
		_, message, err := c.connection.ReadMessage()

		if err != nil {
			if err.Error() == "websocket: close 1001 " {
				log.Printf("Client from %s disconnected", c.GetRemoteAddr())
			} else {
				log.Printf("Client from %s error: %s", c.GetRemoteAddr(), err.Error())
			}

			// Tell the routine to stop
			c.stop <- true
			break
		} else {
			select {
			case c.incoming <- message:
				// Message was sent to incoming queue
			default:
				// Incoming queue was full
				log.Printf("Client from %s filled incoming message queue, disconnecting.", c.GetRemoteAddr())
				c.stop <- true
				return
			}
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
	c.connection = conn
	c.request = req
	c.permissions = auth.Permissions{}
	c.connected = true
	c.subscriptions = []string{}
	c.write = permissionCache{}
	c.stopping = false
	c.outgoing = make(chan []byte, OUTGOING_BUFFER)
	c.incoming = make(chan []byte, INCOMING_BUFFER)
	c.stop = make(chan bool)
	c.mutex = &sync.Mutex{}

	connectedClients++

	return &c
}
