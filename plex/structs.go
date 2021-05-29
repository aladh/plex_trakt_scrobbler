package plex

import (
	"regexp"
)

const ShowType = "show"
const MovieType = "movie"

type Payload struct {
	Account  Account
	Event    string
	Metadata Metadata
	Server   Server
}

type Metadata struct {
	IDs []struct {
		URI string `json:"id"`
	} `json:"Guid"`
	LibrarySectionType string

	// This isn't used but the JSON parsing doesn't work without it since it's case insensitive
	// https://github.com/golang/go/issues/14750
	GUID string `json:"guid"`
}

type Server struct {
	UUID string `json:"uuid"`
}

type Account struct {
	Title string `json:"title"`
}

var idURIRegex = regexp.MustCompile(`(\w*)://(\w*)`)

func (p *Payload) Type() string {
	return p.Metadata.LibrarySectionType
}

func (p *Payload) IDs() map[string]string {
	ids := map[string]string{}

	for _, id := range p.Metadata.IDs {
		matches := idURIRegex.FindStringSubmatch(id.URI)
		ids[matches[1]] = matches[2]
	}

	return ids
}
