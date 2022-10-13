package config

import (
	"fmt"
	"strings"
)

type Profiles []*Profile

type Profile struct {
	Name    string `json:"name"`
	Alias   string `json:"alias"`
	URL     string `json:"url"`
	Token   string `json:"token"`
	Current bool   `json:"current"`
}

func (profile *Profile) Redact() Profile {
	if len(profile.Token) > 10 {
		profile.Token = profile.Token[0:9] + "XXX..."
	}
	return *profile
}

func (profile *Profile) IsValid() (bool, string) {
	if len(profile.URL) < 10 || !strings.HasPrefix(profile.URL, "https://") {
		return false, fmt.Sprintf("URL %s doesn't seem to be valid. Check your profile via the `gabi config currentprofile` command \n", profile.URL)
	}

	if len(profile.Token) < 10 {
		return false, fmt.Sprintf("Token doesn't seem to be valid. Check your profile via the `gabi config currentprofile` command \n")
	}

	return true, ""
}
