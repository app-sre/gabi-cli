package gabi

import (
	"net/http"
	"strings"

	"github.com/cristianoveiga/gabi-cli/pkg/config"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	Token      string
}

func NewClient(profile config.Profile) (*Client, error) {
	c := Client{
		httpClient: &http.Client{},
		baseURL:    strings.TrimSuffix(profile.URL, "/"),
		Token:      profile.Token,
	}
	return &c, nil
}
