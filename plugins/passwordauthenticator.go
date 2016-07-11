package plugins

import (
	"github.com/PuSuEngine/pusud/auth"
	"github.com/PuSuEngine/pusud/core"
	"gopkg.in/yaml.v2"
	"log"
)

type passwordAuthenticator struct {
}

type passwordSettings struct {
	Passwords map[string][]string	`yaml:"passwords"`
}

var settings *passwordSettings
var DEBUG = false

func getSettings() *passwordSettings {
	if settings == nil {
		data := core.GetSettingsContents()
		s := passwordSettings{}
		yaml.Unmarshal(*data, &s)
		settings = &s
	}

	return settings
}

func (ma passwordAuthenticator) GetPermissions(authorization string) auth.Permissions {
	d := auth.Permissions{}

	if authorization == "" {
		if DEBUG {
			log.Printf("No password provided")
		}
	} else {
		match := false
		s := getSettings()
		// map is password -> list of channels
		for k, v := range s.Passwords {
			if k == authorization {
				for _, c := range v {
					d[c] = &auth.Permission{true, true}
					if DEBUG {
						log.Printf("Password %s gave access to %s", k, c)
					}
					match = true
				}
			}
		}

		if !match {
			if DEBUG {
				log.Printf("Invalid password, got no access to anything")
			}
		}
	}

	// Everyone has READ to all for now
	if _, ok := d["*"]; !ok {
		d["*"] = &auth.Permission{true, false}
	}


	return d
}

func init() {
	pa := passwordAuthenticator{}

	auth.RegisterAuthenticator("PasswordAuthenticator", pa)
}
