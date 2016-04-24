package core

import (
	"fmt"
	"log"
	"github.com/lietu/pusud/auth"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func StartListeners(settings *Settings, authenticator auth.Authenticator) {
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

	c := NewClient(conn, r)
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

func AllowOrigin(r *http.Request) bool {
	return true
}

func init() {
	upgrader.CheckOrigin = AllowOrigin
}