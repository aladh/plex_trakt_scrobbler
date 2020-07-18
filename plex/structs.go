package plex

import (
	"regexp"
)

const ShowType = "show"
const movieType = "movie"

type Payload struct {
	Event    string
	Metadata Metadata
	Server   Server
}

type Metadata struct {
	GrandparentGUID    string
	GrandparentTitle   string
	GUID               string
	Index              int
	LibrarySectionType string
	ParentIndex        int
}

type Server struct {
	UUID string `json:"uuid"`
}

type ID struct {
	Provider string
	Value    string
}

var idRegex = regexp.MustCompile(`.*://(.*)\?`)

func (p *Payload) Type() string {
	return p.Metadata.LibrarySectionType
}

func (m *Metadata) Title() string {
	return m.GrandparentTitle
}

func (m *Metadata) Season() int {
	return m.ParentIndex
}

func (m *Metadata) Episode() int {
	return m.Index
}

func (m *Metadata) ID() string {
	guid := m.GrandparentGUID

	if m.LibrarySectionType == movieType {
		guid = m.GUID
	}

	matches := idRegex.FindStringSubmatch(guid)

	return matches[1]
}
