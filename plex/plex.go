package plex

import (
	"fmt"
	"regexp"
	"strconv"
)

type Payload struct {
	Event    string
	Metadata Metadata
	Server   Server
}

type Metadata struct {
	LibrarySectionType string
	GrandparentTitle   string
	ParentIndex        int
	Index              int
	GrandparentGUID    string
}

type Server struct {
	UUID string `json:"uuid"`
}

type ID struct {
	Provider string
	Value    int
}

var idRegex = regexp.MustCompile(`.*\.(.*)://(.*)\?`)

func (m *Metadata) Title() string {
	return m.GrandparentTitle
}

func (m *Metadata) Season() int {
	return m.ParentIndex
}

func (m *Metadata) Episode() int {
	return m.Index
}

func (m *Metadata) ID() (*ID, error) {
	matches := idRegex.FindStringSubmatch(m.GrandparentGUID)

	value, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("error converting to integer: %w", err)
	}

	return &ID{Provider: matches[1], Value: value}, nil
}
