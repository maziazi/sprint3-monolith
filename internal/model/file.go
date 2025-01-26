package model

type File struct {
	ID           int    `json:"id"`
	URI          string `json:"uri"`
	ThumbnailURI string `json:"thumbnailUri"`
}
