<p align="center">
  <h3 align="center">Plex Trakt Scrobbler</h3>

  <p align="center">
    Scrobbles movies and TV shows that you watch on Plex 
</p>

## About The Project

This project leverages Plex Media Server's [webhooks](https://support.plex.tv/articles/115002267687-webhooks/) feature,
and runs a server that scrobbles movies and TV shows you watch to Trakt.

### Prerequisites

- Trakt application
- Plex Media Server (Plex Pass required for webhooks)

## Usage

A Docker image can be built using the provided Dockerfile. The environment variables required to run the application
are:

| Environment variable       | Description                                             |
|----------------------------|---------------------------------------------------------|
| TRAKT_CLIENT_ID            | The client ID of your Trakt application                 |
| TRAKT_ACCESS_TOKEN         | The access token provided to your application via OAuth |
| PLEX_SERVER_UUIDS          | A pipe-separated list of server UUIDs to allow          |
| PLEX_USERNAME              | The username of the plex account to scrobble for        |
| PORT                       | The port to listen on (defaults to 8080)                |
| ERROR_WEBHOOK_URL          | A URL where error notifications are sent                |
| MOVIE_SCROBBLE_WEBHOOK_URL | A URL where movie scrobble notifications are sent       |
