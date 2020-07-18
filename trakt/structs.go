package trakt

import "time"

const tvdb = "tvdb"

type WatchEpisodeRequest struct {
	Shows []Show `json:"shows"`
}

type Show struct {
	Ids     map[string]string `json:"ids"`
	Seasons []Season          `json:"seasons"`
}

type Season struct {
	Number   int       `json:"number"`
	Episodes []Episode `json:"episodes"`
}

type Episode struct {
	Number    int    `json:"number"`
	WatchedAt string `json:"watched_at"`
}

func watchEpisodeRequest(id string, season, episode int, watchedAt time.Time) *WatchEpisodeRequest {
	return &WatchEpisodeRequest{
		Shows: []Show{
			{
				Ids: map[string]string{tvdb: id},
				Seasons: []Season{
					{
						Number: season,
						Episodes: []Episode{
							{
								Number:    episode,
								WatchedAt: watchedAt.Format(time.RFC3339),
							},
						},
					},
				},
			},
		},
	}
}
