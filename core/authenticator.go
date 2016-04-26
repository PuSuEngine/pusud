package core

import (
	"github.com/PuSuEngine/pusud/auth"
)

var authenticator auth.Authenticator

// Set the active Authenticator, should probably only be called once, from main()
func SetAuthenticator(a auth.Authenticator) {
	authenticator = a
}

func getAuthenticator() auth.Authenticator {
	return authenticator
}
