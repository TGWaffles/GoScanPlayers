package webhook

import (
	"GoScanPlayers/models"
	playerLib "GoScanPlayers/player"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type PostData struct {
	Username  string         `json:"username"`
	Content   string         `json:"content"`
	Embeds    []models.Embed `json:"embeds"`
	AvatarUrl string         `json:"avatar_url"`
}

func PostDataToURL(webhookUrl string, webhookContent string, player *models.Player) {
	username := playerLib.LookupSingularUUID(player.Uuid)
	onlineText := "ONLINE"
	colour := 0x00FF00
	if !player.IsOnline {
		onlineText = "OFFLINE"
		colour = 0xFF0000
	}
	var fields []models.Field
	if player.Note != "" {
		fields = []models.Field{
			{
				Name:   "Note",
				Value:  player.Note,
				Inline: false,
			},
		}
	}
	headUrl := fmt.Sprintf("https://cravatar.eu/helmavatar/%s/256.png", player.Uuid)
	data := PostData{
		Username: "Player Notifier",
		Content:  webhookContent,
		Embeds: []models.Embed{{
			Title:       fmt.Sprintf("%s: %s", username, onlineText),
			Description: fmt.Sprintf("%s status: %s", username, player.OnlineStatus),
			Author: models.Author{
				Name:    username,
				IconUrl: headUrl,
				Url:     fmt.Sprintf("https://sky.shiiyu.moe/stats/%s", player.Uuid),
			},
			Colour: colour,
			Fields: fields,
		}},
		AvatarUrl: headUrl,
	}
	postData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	_, _ = http.Post(webhookUrl, "application/json", bytes.NewReader(postData))
}
