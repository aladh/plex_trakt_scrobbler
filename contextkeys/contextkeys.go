package contextkeys

type key string

const (
	Config      key = "config"
	TraktClient key = "trakt_client"
	Request     key = "request"
	Payload     key = "payload"
)
