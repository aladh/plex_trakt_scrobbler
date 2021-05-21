package plex

import (
	"regexp"
)

const ShowType = "show"
const MovieType = "movie"

type Payload struct {
	Event    string
	Metadata Metadata
	Server   Server
}

type Metadata struct {
	IDs []struct {
		Uri string `json:"id"`
	} `json:"Guid"`
	LibrarySectionType string

	// This isn't used but the JSON parsing doesn't work without it since it's case insensitive
	// https://github.com/golang/go/issues/14750
	Guid string `json:"guid"`
}

type Server struct {
	UUID string `json:"uuid"`
}

var idUriRegex = regexp.MustCompile(`(\w*)://(\w*)`)

func (p *Payload) Type() string {
	return p.Metadata.LibrarySectionType
}

func (p *Payload) IDs() map[string]string {
	ids := map[string]string{}

	for _, id := range p.Metadata.IDs {
		matches := idUriRegex.FindStringSubmatch(id.Uri)
		ids[matches[1]] = matches[2]
	}

	return ids
}
