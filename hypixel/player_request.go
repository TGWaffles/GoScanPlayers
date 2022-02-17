package hypixel

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"tests/storage"
	"time"
)

const ApiUrl = "https://api.hypixel.net"

type RequestMaker struct {
	isRateLimited bool
	keyIsValid         bool
	rateLimitResetTime time.Time
	data *storage.Data
}

func New(data *storage.Data) *RequestMaker {
	handler := &RequestMaker{
		isRateLimited:      false,
		keyIsValid:         true,
		rateLimitResetTime: time.Now(),
		data:               data,
	}
	handler.updatePlayersLoop()
	return handler
}

func (handler *RequestMaker) updatePlayersLoop() {

}

type StatusRequestResponse struct {
	Success bool `json:"success"`
	Cause string `json:"cause"`
	Uuid string `json:"uuid"`
	Session SessionResponse `json:"session"`
}

type SessionResponse struct {
	Online bool `json:"online"`
	GameType string `json:"gameType"`
	Mode string `json:"mode"`
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
