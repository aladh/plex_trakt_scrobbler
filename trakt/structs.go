package trakt

import "time"

type WatchEpisodeRequest struct {
	Episodes []Episode `json:"episodes"`
}

type Episode struct {
	Ids       map[string]string `json:"ids"`
	WatchedAt string            `json:"watched_at"`
}

type WatchMovieRequest struct {
	Movies []Movie `json:"movies"`
}

type Movie struct {
	Ids       map[string]string `json:"ids"`
	WatchedAt string            `json:"watched_at"`
}

func watchEpisodeRequest(ids map[string]string, watchedAt time.Time) *WatchEpisodeRequest {
	return &WatchEpisodeRequest{
		Episodes: []Episode{
			{
				Ids:       ids,
				WatchedAt: watchedAt.Format(time.RFC3339),
			},
		},
	}
}

func watchMovieRequest(ids map[string]string, watchedAt time.Time) *WatchMovieRequest {
	return &WatchMovieRequest{
		Movies: []Movie{
			{
				Ids:       ids,
				WatchedAt: watchedAt.Format(time.RFC3339),
			},
		},
	}
}
