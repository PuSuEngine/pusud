package auth

type NoAuthenticator struct {

}

func (na NoAuthenticator) GetPermissions(authorization string) Permissions {
	d := Permissions{};
	d["*"] = &Permission{true, true}
	return d
}

func init() {
	na := NoAuthenticator{}
	RegisterAuthenticator("None", na)
}

