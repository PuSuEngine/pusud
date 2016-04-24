package auth

import (
	"log"
)

var authenticators map[string]Authenticator = map[string]Authenticator{}

type Authenticator interface {
	GetPermissions(authorization string) Permissions
}

func RegisterAuthenticator(name string, auth Authenticator) bool {
	_, ok := authenticators[name]

	if ok {
		log.Printf("Trying to re-register Authenticator \"%s\".", name)
		log.Printf("")
		log.Printf("Check your plugins.")
		return false
	}

	authenticators[name] = auth

	return true
}

func GetAuthenticator(name string) (Authenticator, bool) {
	auth, ok := authenticators[name]

	if !ok {
		log.Printf("Authenticator \"%s\" does not exist.", name)
		log.Printf("")

		log.Printf("Valid options are:")
		for k, _ := range authenticators {
			log.Printf("%s", k)
		}

		log.Printf("")
		log.Printf("Check your settings.")

		return nil, false
	}

	return auth, true
}

func GetChannelPermissions(name string, permissions Permissions) (bool, bool) {
	var read bool
	var write bool

	for k, p := range permissions {
		if ChannelMatch(name, k) {
			if p.Read {
				read = true
			}
			if p.Write {
				write = true
			}
		}
	}

	return read, write
}

func ChannelMatch(channel string, match string) bool {
	if match == "*" {
		return true
	}

	// Figure out how many characters of the channel must match (wildcard)
	var minMatch int
	for _, c := range match {
		if c == '*' {
			break
		}

		minMatch++
	}

	if minMatch <= len(channel) {
		log.Printf("%s vs. %s", channel[:minMatch], match[:minMatch])

		if channel[:minMatch] == match[:minMatch] {
			return true
		}
	}

	return false
}
