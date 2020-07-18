package trakt

import "time"

const tvdb = "tvdb"
const imdb = "imdb"

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

type WatchMovieRequest struct {
	Movies []Movie `json:"movies"`
}

type Movie struct {
	Ids       map[string]string `json:"ids"`
	WatchedAt string            `json:"watched_at"`
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

func watchMovieRequest(id string, watchedAt time.Time) *WatchMovieRequest {
	return &WatchMovieRequest{
		Movies: []Movie{
			{
				Ids:       map[string]string{imdb: id},
				WatchedAt: watchedAt.Format(time.RFC3339),
			},
		},
	}
}
