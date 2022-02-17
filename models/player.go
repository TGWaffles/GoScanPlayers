package models

import "fyne.io/fyne/v2/widget"

type Player struct {
	Uuid string `json:"uuid"`
	Note string `json:"note"`
	IsOnline bool `json:"isOnline"`
	HasApiEnabled bool `json:"useSkyblockAPI"`
	OnlineLabel *widget.Label
}


