package auth

import (
	"log"
)

var authenticators map[string]Authenticator = map[string]Authenticator{}

// An Authenticator figures out what permissions the user should have based on
// the authorization string they send.
type Authenticator interface {
	GetPermissions(authorization string) Permissions
}

// Register a new Authenticator. You need to call this in your plugin for the
// Authenticator to be available for use.
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

// Get an authenticator by the name it was registered. Not very useful outside
// of main() function.
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

// Figure out what permissions the given list of permissions gives to the named
// channel. Merges permissions so if you have e.g. read on "*" and write on "foo"
// You will have both read & write on "foo".
func GetChannelPermissions(name string, permissions Permissions) (read bool, write bool) {
	read = false
	write = false

	for k, p := range permissions {
		if channelMatch(name, k) {
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

func channelMatch(channel string, match string) bool {
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
