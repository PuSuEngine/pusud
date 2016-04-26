package plugins

import "github.com/PuSuEngine/pusud/auth"

type myAuthenticator struct {
}

func (ma myAuthenticator) GetPermissions(authorization string) auth.Permissions {
	d := auth.Permissions{}
	return d
}

func init() {
	ma := myAuthenticator{}

	auth.RegisterAuthenticator("MyAuthenticator", ma)
}
