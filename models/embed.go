package models

type Embed struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Url         string    `json:"url"`
	Image       Image     `json:"image"`
	Thumbnail   Thumbnail `json:"thumbnail"`
	Author      Author    `json:"author"`
	Fields      []Field   `json:"fields"`
	Colour      int       `json:"color"`
}

type Image struct {
	Url string `json:"url"`
}

type Thumbnail struct {
	Url string `json:"url"`
}

type Author struct {
	Name    string `json:"name"`
	Url     string `json:"url"`
	IconUrl string `json:"icon_url"`
}

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}
