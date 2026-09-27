package games

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/kartFr/Asset-Reuploader/internal/retry"
	"github.com/kartFr/Asset-Reuploader/internal/roblox"
)

func newGamesHandler(c *roblox.Client, url string) (func() (*GamesResponse, error), error) {
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return nil, err
	}

	return func() (*GamesResponse, error) {
		req.AddCookie(&http.Cookie{
			Name:  ".ROBLOSECURITY",
			Value: c.Cookie,
		})

		resp, err := c.DoRequest(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, errors.New(resp.Status)
		}

		var gamesResponse GamesResponse
		if err := json.NewDecoder(resp.Body).Decode(&gamesResponse); err != nil {
			return nil, err
		}
		return &gamesResponse, nil
	}, nil
}

func NewGroupGamesHandler(c *roblox.Client, groupID int64) (func() (*GamesResponse, error), error) {
	url := fmt.Sprintf("https://games.roblox.com/v2/groups/%d/gamesV2?limit=100", groupID)
	return newGamesHandler(c, url)
}

func GroupGames(c *roblox.Client, groupID int64) (*GamesResponse, error) {
	handler, err := NewGroupGamesHandler(c, groupID)
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
