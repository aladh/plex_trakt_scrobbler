package plex

import (
	"regexp"
)

const showType = "show"

type Payload struct {
	Event    string
	Metadata Metadata
	Server   Server
}

type Metadata struct {
	GrandparentGUID    string
	GrandparentTitle   string
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

	matches := idRegex.FindStringSubmatch(guid)

	return matches[1]
}
