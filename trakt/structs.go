package trakt

import "time"

type WatchEpisodeRequest struct {
	Shows []Show `json:"shows"`
}

type Show struct {
	Ids     IDs      `json:"ids"`
	Seasons []Season `json:"seasons"`
}

type IDs struct {
	Tvdb string `json:"tvdb"`
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
				Ids: IDs{Tvdb: id},
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
