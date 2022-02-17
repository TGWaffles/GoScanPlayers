package player

import (
	"GoScanPlayers/models"
	"GoScanPlayers/storage"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ListHandler struct {
	index  int
	data   *storage.Data
	window fyne.Window
	List   *widget.List
	uuids  map[string]string
	master *fyne.Container
}

func (handler *ListHandler) getNextItem() fyne.CanvasObject {
	if handler.index >= len(handler.data.Players) || len(handler.data.Players) == 0 {
		if len(handler.data.Players) > 0 {
			handler.index = 0
		} else {
			return widget.NewLabel("")
		}
	}
	player := handler.data.Players[handler.index]
	handler.index++
	playerUuid := player.Uuid
	playerName := handler.uuids[playerUuid]
	playerNameLabel := widget.NewLabel(playerName)
	playerNoteLabel := widget.NewLabel(player.Note)
	playerNoteLabel.Alignment = fyne.TextAlignLeading
	nameWithNote := container.NewGridWithColumns(2,
		playerNameLabel,
		playerNoteLabel,
	)
	onlineStatus := widget.NewLabel(player.OnlineStatus)
	onlineStatus.Alignment = fyne.TextAlignCenter
	player.OnlineLabel = onlineStatus
	return container.NewBorder(
		nil,
		nil,
		nameWithNote,
		widget.NewButton("x", func() {
			handler.index = 0
			handler.removePlayer(playerUuid)
			handler.data.Parent.SaveData()
			handler.ReloadList()
		}),
		onlineStatus,
	)
}

func (handler *ListHandler) removePlayer(uuid string) {
	newPlayers := make([]*models.Player, 0)
	for _, player := range handler.data.Players {
		if player.Uuid != uuid {
			newPlayers = append(newPlayers, player)
		}
	}
	handler.data.Players = newPlayers
}

func (handler *ListHandler) refreshUuids() {
	uuids := make([]string, 0)
	for _, player := range handler.data.Players {
		uuids = append(uuids, player.Uuid)
	}
	handler.uuids = LookupUUIDs(uuids).Uuids
}

func GeneratePlayerList(data *storage.Data, window fyne.Window) *ListHandler {
	handler := &ListHandler{
		index:  0,
		data:   data,
		window: window,
	}
	handler.refreshUuids()

	handler.createNewList()

	return handler
}

func (handler *ListHandler) createNewList() *widget.List {
	list := widget.NewList(
		func() int {
			return len(handler.data.Players)
		},
		func() fyne.CanvasObject {
			return handler.getNextItem()
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			return
		},
	)

	list.Resize(fyne.NewSize(500, 150))

	handler.List = list
	return list
}

func (handler *ListHandler) SetMaster(master *fyne.Container) {
	handler.master = master
}

func (handler *ListHandler) AddPlayer(username string, note string) {
	usernames := LookupUsernames([]string{username})
	if usernames.Usernames[username] == "" || usernames.Usernames[username] == "Unknown Player" {
		popUp := widget.NewPopUp(widget.NewLabel("Invalid Username!"), handler.window.Canvas())
		popUp.Show()
		return
	}
	player := &models.Player{
		Uuid:          usernames.Usernames[username],
		Note:          note,
		IsOnline:      false,
		HasApiEnabled: true,
		OnlineStatus:  "NOT CHECKED",
	}
	handler.data.Players = append(handler.data.Players, player)
	handler.refreshUuids()
	handler.ReloadList()
}

func (handler *ListHandler) ReloadList() {
	handler.index = len(handler.data.Players) - 1
	for index, element := range handler.master.Objects {
		if element == handler.List {
			handler.master.Objects[index] = handler.createNewList()
		}
	}
	handler.List.Refresh()
	handler.window.Content().Refresh()
}
