package auth

import (
	"testing"
)

func TestNoAuthenticator(t *testing.T) {
	n, ok := GetAuthenticator("None")

	if !ok {
		t.Fatalf("Failed to load the built-in None authenticator")
	}

	perms := n.GetPermissions("")
	read, write := GetChannelPermissions("", perms)
	if !read || !write {
		t.Fatalf("None authenticator denied permission")
	}
	read, write = GetChannelPermissions("üå123XCZX..`", perms)
	if !read || !write {
		t.Fatalf("None authenticator denied permission")
	}
}