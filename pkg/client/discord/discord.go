package discord

import (
	"encoding/json"
	"fmt"
	. "github.com/gofree-wtf/discord-cloud-status-bot/pkg/common"
	"net/http"
	"time"
)

const DiscordStatusApiUrl = "https://discordstatus.com/api/v2/status.json"

func GetDiscordStatus() (*DiscordStatus, error) {
	result := &DiscordStatus{}

	resp, err := HttpClient.Get(DiscordStatusApiUrl)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("discord status response code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type DiscordStatus struct {
	Page   Page
	Status Status
}

type Page struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	URL        string `json:"url"`
	ITimeZone  string `json:"time_zone"`
	IUpdatedAt string `json:"updated_at"`
}

type Status struct {
	Indicator   string `json:"indicator"`
	Description string `json:"description"`
}

func (p *Page) UpdatedAt() (updatedAt time.Time, err error) {
	var apiLoc *time.Location

	apiLoc, err = time.LoadLocation(p.ITimeZone)
	if err != nil {
		return
	}

	updatedAt, err = time.ParseInLocation(time.RFC3339Nano, p.IUpdatedAt, apiLoc)
	if err != nil {
		return
	}
	return updatedAt.In(Config.Bot.Location), err
}
