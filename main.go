package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"tests/player"
	"tests/storage"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	storageHandler := storage.New()
	data := storageHandler.GetData()
	a := app.New()
	w := a.NewWindow("Player Scanner")

	apiKeyEntry := widget.NewEntry()
	apiKeyEntry.SetPlaceHolder("Enter Your Hypixel Api Key Here!")
	if data.ApiKey != "" {
		apiKeyEntry.Text = data.ApiKey
	}
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

	masterConfig := container.NewVBox(
		widget.NewLabel("Player Scanner: Config"),
		apiKeyEntry,
		webhookUrl,
		webhookMessage,
	)
	playerList := player.GeneratePlayerList(data, w)

	playerNameEntry := widget.NewEntry()
	playerNameEntry.SetPlaceHolder("Player username...")
	playerNoteEntry := widget.NewEntry()
	playerNoteEntry.SetPlaceHolder("Player Note...")
	playerNameConfirm := widget.NewButton("Add Player", func() {
		playerList.AddPlayer(playerNameEntry.Text, playerNoteEntry.Text, w)
		playerNameEntry.Text = ""
		playerNameEntry.Refresh()
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