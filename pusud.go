package main

import (
	"github.com/lietu/pusud/core"
	"github.com/lietu/pusud/auth"
	_ "github.com/lietu/pusud/plugins"
	"log"
)

func main() {
	settings := core.ReadSettings()

	log.Printf("Using %s authenticator", settings.Authenticator)

	authenticator, ok := auth.GetAuthenticator(settings.Authenticator)

	if !ok {
		log.Fatalf("Couldn't find the configured authenticator")
	}

	core.SetAuthenticator(authenticator)
	core.SetupNetwork(settings)
	core.StartListeners(settings, authenticator)
}
