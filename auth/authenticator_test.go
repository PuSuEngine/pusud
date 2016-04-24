package auth

import (
	"testing"
)


type FakeAuthenticator struct {

}

func (na FakeAuthenticator) GetPermissions(name string, authorization string) map[string]Permission {
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
	_, ok := GetAuthenticator("fake")
	if !ok {
		t.Fatalf("Could not load registered authenticator")
	}
	_, ok = GetAuthenticator("fake2")
	if ok {
		t.Fatalf("Loaded a non-existent authenticator")
	}
}
