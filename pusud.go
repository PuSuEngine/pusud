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

	permissions := authenticator.GetPermissions("foo", "bar")

	for k, p := range permissions {
		log.Printf("Channel: %s, Read: %t  Write: %t", k, p.Read, p.Write)
	}

	channels := []string{"user", "user.1234", "users"}

	for _, v := range channels {
		read, write := auth.GetChannelPermissions(v, permissions)
		log.Printf("Channel: %s, Read: %t  Write: %t", v, read, write)
	}
}
