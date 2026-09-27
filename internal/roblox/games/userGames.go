package games

import (
	"fmt"

	"github.com/kartFr/Asset-Reuploader/internal/retry"
	"github.com/kartFr/Asset-Reuploader/internal/roblox"
)

func NewUserGamesHandler(c *roblox.Client, userID int64) (func() (*GamesResponse, error), error) {
	return newGamesHandler(c, fmt.Sprintf("https://games.roblox.com/v2/users/%d/games?limit=50", userID))
}

func UserGames(c *roblox.Client, userID int64) (*GamesResponse, error) {
	handler, err := NewUserGamesHandler(c, userID)
	if err != nil {
		return nil, err
	}

	return retry.Do(
		retry.NewOptions(retry.Tries(3)),
		func(_ int) (*GamesResponse, error) {
			placeDetails, err := handler()
			if err != nil {
				return nil, &retry.ContinueRetry{Err: err}
			}

			return placeDetails, nil
		},
	)
}
