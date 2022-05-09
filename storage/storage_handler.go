package storage

import (
	"GoScanPlayers/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type Handler struct {
	fileName string
	data     *Data
}

type Data struct {
	UpdateFrequency      int  `json:"updateFrequency"`
	EnableSkyblockLookup bool `json:"enableSkyblock"`

	ApiKey         string `json:"apiKey"`
	WebhookUrl     string `json:"webhookUrl"`
	WebhookContent string `json:"webhookContent"`

	Players []*models.Player `json:"players"`
	Parent  *Handler         `json:"-"`
}

func New() *Handler {
	handler := Handler{fileName: "data.json", data: nil}
	return &handler
}

func (handler *Handler) GetData() *Data {
	if handler.data != nil {
		return handler.data
	}
	handler.data = handler.loadData()
	return handler.data
}

func (handler *Handler) loadData() *Data {
	byteData, err := ioutil.ReadFile(handler.fileName)
	if err != nil {
		if strings.Contains(err.Error(), "The system cannot find the file specified.") {
			handler.createFile()
			return handler.loadData()
		}
		log.Fatal(err)
	}
	data := &Data{Parent: handler}
	err = json.Unmarshal(byteData, data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func (handler *Handler) createFile() {
	handler.data = &Data{
		UpdateFrequency:      60,
		EnableSkyblockLookup: true,
		WebhookUrl:           "",
		WebhookContent:       "",
		Players:              make([]*models.Player, 0),
		Parent:               handler,
	}
	handler.SaveData()
}

func (handler *Handler) SaveData() {
	data, err := json.Marshal(handler.data)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(handler.fileName, data, 0755)
	if err != nil {
		fmt.Printf("Unable to write file: %v", err)
	}
}
