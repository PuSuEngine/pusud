package main

import (
	"github.com/PuSuEngine/pusud/auth"
	"github.com/PuSuEngine/pusud/core"
	_ "github.com/PuSuEngine/pusud/plugins"
	"log"
)

func main() {
	settings := core.GetSettings()

	log.Printf("Using %s authenticator", settings.Authenticator)

	authenticator, ok := auth.GetAuthenticator(settings.Authenticator)

	if !ok {
		log.Fatalf("Couldn't find the configured authenticator")
	}

	core.SetAuthenticator(authenticator)
	core.SetupNetwork(settings)
	core.StartListeners(settings, authenticator)
}
