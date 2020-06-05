package plex

import (
	"log"
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
	GrandparentGuid    string
}

type Server struct {
	UUID string `json:"uuid"`
}

type ID struct {
	Provider string
	Value    int
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

func (m *Metadata) ID() ID {
	re := regexp.MustCompile(".*\\.(.*)://(.*)\\?")
	matches := re.FindStringSubmatch(m.GrandparentGuid)

	value, err := strconv.Atoi(matches[2])
	if err != nil {
		log.Printf("Error getting value from %s\n", m.GrandparentGuid)
	}

	return ID{
		Provider: matches[1],
		Value:    value,
	}
}
