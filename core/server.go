package core

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/PuSuEngine/pusud/auth"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{}

func StartListeners(settings *Settings, authenticator auth.Authenticator) {
	go statusMonitor()
	go runNetworkListener(settings.NetworkPort)
	runClientListener(settings.ClientPort, authenticator)
}

func runNetworkListener(port int) {
	// TODO: Listen for relay<->relay connections
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	log.Printf("Connection from %s", r.RemoteAddr)

	if err != nil {
		log.Print("upgrade: ", err)
		return
	}

	c := newClient(conn, r)
	defer c.Close()

	c.SendHello()
	c.Handle()
}

func runClientListener(port int, authenticator auth.Authenticator) {
	address := fmt.Sprintf("0.0.0.0:%d", port)

	log.Printf("Starting to listen for client connections on %s", address)

	http.HandleFunc("/", websocketHandler)

	log.Fatal(http.ListenAndServe(address, nil))
}

func statusMonitor() {
	for {
		time.Sleep(time.Second * 30)
		log.Printf("")
		log.Printf("----- Status update -----")
		log.Printf("Currently have %d connected clients", connectedClients)

		log.Printf("Received: %d messages", readCounter)
		readCounter = 0

		log.Printf("Sent    : %d messages", writeCounter)
		writeCounter = 0
		log.Printf("")
	}
}

func allowOrigin(r *http.Request) bool {
	return true
}

func init() {
	upgrader.CheckOrigin = allowOrigin
}
