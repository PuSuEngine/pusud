package plugins

import "github.com/lietu/pusud/auth"

type MyAuthenticator struct {

}

func (ma MyAuthenticator) GetPermissions(authorization string) map[string]auth.Permission {
	d := map[string]auth.Permission{}
	d["*"] = auth.Permission{true, false}
	d["user.*"] = auth.Permission{true, true}
	return d
}

func init() {
	ma := MyAuthenticator{}

	auth.RegisterAuthenticator("MyAuthenticator", ma)
}
