package core

import (
	"github.com/lietu/pusud/auth"
)

var authenticator auth.Authenticator

func SetAuthenticator(a auth.Authenticator) {
	authenticator = a
}

func GetAuthenticator() auth.Authenticator {
	return authenticator
}
