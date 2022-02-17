package main

import (
	"GoScanPlayers/hypixel"
	"GoScanPlayers/player"
	"GoScanPlayers/storage"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	storageHandler := storage.New()
	data := storageHandler.GetData()
	a := app.New()
	w := a.NewWindow("Player Scanner")
	playerList := player.GeneratePlayerList(data, w)
	apiKeyEntry := widget.NewEntry()
	apiKeyEntry.SetPlaceHolder("Enter Your Hypixel Api Key Here!")
	if data.ApiKey != "" {
		apiKeyEntry.Text = data.ApiKey
	}
	hypixelLookupChecker := hypixel.New(data, playerList)
	webhookMessage := widget.NewEntry()
	webhookMessage.SetPlaceHolder("Discord Webhook Message Content")
	if data.WebhookContent != "" {
		webhookMessage.Text = data.WebhookContent
	}
	webhookUrl := widget.NewEntry()
	webhookUrl.SetPlaceHolder("Discord Webhook URL")
	if data.WebhookUrl != "" {
		webhookUrl.Text = data.WebhookUrl
	}

	saveConfigButton := widget.NewButton("Save Config", func() {
		data.WebhookUrl = webhookUrl.Text
		data.WebhookContent = webhookMessage.Text
		data.ApiKey = apiKeyEntry.Text
		storageHandler.SaveData()
		hypixelLookupChecker.ApiKeyUpdated()
	})

	masterConfig := container.NewVBox(
		widget.NewLabel("Player Scanner: Config"),
		apiKeyEntry,
		webhookUrl,
		webhookMessage,
		saveConfigButton,
	)

	playerNameEntry := widget.NewEntry()
	playerNameEntry.SetPlaceHolder("Player username...")
	playerNoteEntry := widget.NewEntry()
	playerNoteEntry.SetPlaceHolder("Player Note...")
	playerNameConfirm := widget.NewButton("Add Player", func() {
		go playerList.AddPlayer(playerNameEntry.Text, playerNoteEntry.Text)
		playerNameEntry.Text = ""
		playerNoteEntry.Text = ""
		playerNameEntry.Refresh()
		playerNoteEntry.Refresh()
		storageHandler.SaveData()
	})

	playerInput := container.NewVBox(
		playerNameEntry,
		playerNoteEntry,
		playerNameConfirm,
	)

	master := container.NewGridWithRows(3,
		masterConfig,
		playerInput,
		playerList.List,
	)

	playerList.SetMaster(master)

	w.SetContent(master)

	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}
