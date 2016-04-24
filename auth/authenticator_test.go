package auth

import (
	"testing"
)


type FakeAuthenticator struct {

}

func (na FakeAuthenticator) GetPermissions(authorization string) map[string]Permission {
	d := map[string]Permission{};
	return d
}


func TestRegisterAuthenticator(t *testing.T) {
	if !RegisterAuthenticator("fake", FakeAuthenticator{}) {
		t.Fatalf("Failed to register authenticator")
	}
	if RegisterAuthenticator("fake", FakeAuthenticator{}) {
		t.Fatalf("Re-registering authenticator caused no error")
	}
}

func TestGetAuthenticator(t *testing.T) {
	// Ignore errors, they're not important
	RegisterAuthenticator("fake", FakeAuthenticator{})

	_, ok := GetAuthenticator("fake")
	if !ok {
		t.Fatalf("Could not load registered authenticator")
	}
	_, ok = GetAuthenticator("fake2")
	if ok {
		t.Fatalf("Loaded a non-existent authenticator")
	}
}

func TestGetChannelPermissions(t *testing.T) {
	perms := map[string]Permission{}
	perms["user.*"] = Permission{true, false}
	perms["write.*"] = Permission{false, true}

	read, write := GetChannelPermissions("users", perms)
	if read || write {
		t.Fatalf("Got invalid permissions")
	}

	read, write = GetChannelPermissions("user", perms)
	if read || write {
		t.Fatalf("Got invalid permissions")
	}

	read, write = GetChannelPermissions("user.1", perms)
	if !read || write {
		t.Fatalf("Got invalid permissions")
	}

	read, write = GetChannelPermissions("write.1", perms)
	if read || !write {
		t.Fatalf("Got invalid permissions")
	}

	// Give read access to everything
	perms["*"] = Permission{true, false}
	// Now "write.*" should have both read and write

	read, write = GetChannelPermissions("write.1", perms)
	if !read || !write {
		t.Fatalf("Got invalid permissions")
	}
}
