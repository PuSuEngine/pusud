package auth

type Permission struct {
	Read  bool
	Write bool
}

type Permissions map[string]*Permission
