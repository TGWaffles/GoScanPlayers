package player

import (
	"GoScanPlayers/models"
	"GoScanPlayers/storage"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type ListHandler struct {
	data   *storage.Data
	window fyne.Window
	List   *widget.List
	uuids  map[string]string
	master *fyne.Container
}

func (handler *ListHandler) createBlankItem() fyne.CanvasObject {
	return container.NewBorder(nil, nil, widget.NewLabel("Username"), widget.NewButton("x", func() {}), widget.NewLabel("OFFLINE"))
}

func (handler *ListHandler) updateItem(index int, object fyne.CanvasObject) {
	player := handler.data.Players[index]
	playerUuid := player.Uuid
	playerName := handler.uuids[playerUuid]
	playerNameLabel := widget.NewLabel(playerName)
	playerNoteLabel := widget.NewLabel(player.Note)
	playerNoteLabel.Alignment = fyne.TextAlignLeading
	nameWithNote := container.NewGridWithColumns(2,
		playerNameLabel,
		playerNoteLabel,
	)

	deleteButton := widget.NewButton("x", func() {
		handler.removePlayer(playerUuid)
		handler.data.Parent.SaveData()
		handler.ReloadList()
	})
	newLayout := layout.NewBorderLayout(
		nil,
		nil,
		nameWithNote,
		deleteButton,
	)
	playerContainer := object.(*fyne.Container)
	playerContainer.Layout = newLayout
	playerContainer.Objects = []fyne.CanvasObject{player.OnlineLabel, nameWithNote, deleteButton}
	playerContainer.Layout.Layout(playerContainer.Objects, playerContainer.Size())
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
		data:   data,
		window: window,
	}
	handler.refreshUuids()

	handler.createNewList()

	return handler
}

func (handler *ListHandler) createNewList() *widget.List {
	for _, player := range handler.data.Players {
		onlineStatus := widget.NewLabel(player.OnlineStatus)
		onlineStatus.Alignment = fyne.TextAlignCenter
		player.OnlineLabel = onlineStatus
	}
	list := widget.NewList(
		func() int {
			return len(handler.data.Players)
		},
		func() fyne.CanvasObject {
			return handler.createBlankItem()
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			handler.updateItem(id, object)
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
	onlineStatus := widget.NewLabel("NOT CHECKED")
	onlineStatus.Alignment = fyne.TextAlignCenter
	player := &models.Player{
		Uuid:          usernames.Usernames[username],
		Note:          note,
		IsOnline:      false,
		HasApiEnabled: true,
		OnlineStatus:  "NOT CHECKED",
		OnlineLabel:   onlineStatus,
	}
	handler.data.Players = append(handler.data.Players, player)
	handler.refreshUuids()
	handler.ReloadList()
}

func (handler *ListHandler) ReloadList() {
	handler.List.Refresh()
	handler.window.Content().Refresh()
}
