package trakt

import "time"

type WatchEpisodeRequest struct {
	Episodes []Episode `json:"episodes"`
}

type Episode struct {
	IDs       map[string]string `json:"ids"`
	WatchedAt string            `json:"watched_at"`
}

type WatchMovieRequest struct {
	Movies []Movie `json:"movies"`
}

type Movie struct {
	IDs       map[string]string `json:"ids"`
	WatchedAt string            `json:"watched_at"`
}

type WatchedMovie struct {
	ID    int
	Movie struct {
		Title string
	}
}

type RemoveHistoryRequest struct {
	IDs []int `json:"ids"`
}

func watchEpisodeRequest(ids map[string]string, watchedAt time.Time) WatchEpisodeRequest {
	return WatchEpisodeRequest{
		Episodes: []Episode{
			{
				IDs:       ids,
				WatchedAt: watchedAt.Format(time.RFC3339),
			},
		},
	}
}

func watchMovieRequest(ids map[string]string, watchedAt time.Time) WatchMovieRequest {
	return WatchMovieRequest{
		Movies: []Movie{
			{
				IDs:       ids,
				WatchedAt: watchedAt.Format(time.RFC3339),
			},
		},
	}
}

func removeHistoryRequest(id int) RemoveHistoryRequest {
	return RemoveHistoryRequest{
		IDs: []int{id},
	}
}
