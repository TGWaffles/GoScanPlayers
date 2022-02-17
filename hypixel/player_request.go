package hypixel

import (
	"GoScanPlayers/player"
	"GoScanPlayers/storage"
	"GoScanPlayers/webhook"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const ApiUrl = "https://api.hypixel.net"

type RequestMaker struct {
	isRateLimited      bool
	keyIsValid         bool
	rateLimitResetTime time.Time
	data               *storage.Data
	playerIndex        int
	playerList         *player.ListHandler
}

func New(data *storage.Data, list *player.ListHandler) *RequestMaker {
	handler := &RequestMaker{
		isRateLimited:      false,
		keyIsValid:         true,
		rateLimitResetTime: time.Now(),
		data:               data,
		playerIndex:        0,
		playerList:         list,
	}
	go handler.updatePlayersLoop()
	return handler
}

func (handler *RequestMaker) getSleepTime() time.Duration {
	// Refresh every .6s, taking 100 req/min
	return 600 * time.Millisecond
}

func (handler *RequestMaker) updatePlayersLoop() {
	time.Sleep(5 * time.Second)
	for {
		for !handler.keyIsValid || len(handler.data.Players) == 0 {
			time.Sleep(5 * time.Second)
		}
		playerObject := handler.data.Players[handler.playerIndex]
		loginText := handler.CheckPlayerOnline(playerObject.Uuid, handler.data.ApiKey)
		if playerObject.LastSuccessfulStatus != loginText {
			playerObject.OnlineStatus = loginText
			playerObject.IsOnline = loginText[:6] == "ONLINE"
			if loginText != "RATE LIMITED" {
				go webhook.PostDataToURL(handler.data.WebhookUrl, handler.data.WebhookContent, playerObject)
			}

			if loginText[:6] == "ONLINE" || loginText[:7] == "OFFLINE" {
				playerObject.LastSuccessfulStatus = loginText
			}
		}
		handler.data.Parent.SaveData()
		playerObject.OnlineLabel.Text = loginText
		playerObject.OnlineLabel.Refresh()
		handler.playerList.ReloadList()
		time.Sleep(handler.getSleepTime())
		handler.playerIndex++
		if handler.playerIndex >= len(handler.data.Players) {
			handler.playerIndex = 0
		}
	}
}

type StatusRequestResponse struct {
	Success bool            `json:"success"`
	Cause   string          `json:"cause"`
	Uuid    string          `json:"uuid"`
	Session SessionResponse `json:"session"`
}

type SessionResponse struct {
	Online   bool   `json:"online"`
	GameType string `json:"gameType"`
	Mode     string `json:"mode"`
}

func (handler *RequestMaker) CheckPlayerOnline(uuid string, apiKey string) string {
	resp, err := http.Get(ApiUrl + "/status?uuid=" + uuid + "&key=" + apiKey)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 403 {
		handler.keyIsValid = false
		return "INVALID KEY"
	}
	if resp.StatusCode == 429 {
		handler.isRateLimited = true
		resetTime := resp.Header.Get("RateLimit-Reset")
		if resetTime == "" {
			resetTime = resp.Header.Get("Retry-After")
		}
		estReset, err := strconv.ParseInt(resetTime, 10, 8)
		if err == nil {
			handler.rateLimitResetTime = time.Now().Add(time.Second * time.Duration(estReset))
		}
		return "RATE LIMITED"
	}
	data := &StatusRequestResponse{}
	byteData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(byteData, data)
	if err != nil {
		log.Fatal(err)
	}
	if !data.Success {
		return data.Cause
	}
	if data.Session.Online {
		return "ONLINE - " + data.Session.GameType + " - " + data.Session.Mode
	} else {
		return "OFFLINE"
	}
}

func (handler *RequestMaker) ApiKeyUpdated() {
	handler.keyIsValid = true
}
