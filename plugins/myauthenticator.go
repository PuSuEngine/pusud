package plugins

import "github.com/lietu/pusud/auth"

type MyAuthenticator struct {
}

func (ma MyAuthenticator) GetPermissions(authorization string) auth.Permissions {
	d := auth.Permissions{}
	return d
}

func init() {
	ma := MyAuthenticator{}

	auth.RegisterAuthenticator("MyAuthenticator", ma)
}
