package trakt

import "time"

type WatchShowRequest struct {
	Shows []Show `json:"shows"`
}

type Show struct {
	Ids     IDs      `json:"ids"`
	Seasons []Season `json:"seasons"`
}

type IDs struct {
	Tvdb int `json:"tvdb"`
}

type Season struct {
	Number   int       `json:"number"`
	Episodes []Episode `json:"episodes"`
}

type Episode struct {
	Number    int    `json:"number"`
	WatchedAt string `json:"watched_at"`
}

func newWatchShowRequest(id, season, episode int, watchedAt time.Time) *WatchShowRequest {
	return &WatchShowRequest{
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
