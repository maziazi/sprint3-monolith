package service

import (
	"context"
	"sprint3/internal/model"
	"sprint3/pkg/database"

	"log"
)

func AddFile(fileURL, thumbnailURL string) (*model.File, error) {
	db := database.GetDBPool()
	var file model.File
	err := db.QueryRow(context.Background(), `INSERT INTO file ("fileUri", "fileThumbnailUri") VALUES ($1, $2) RETURNING "fileId", "fileUri", "fileThumbnailUri"`, fileURL, thumbnailURL).Scan(&file.ID, &file.URI, &file.ThumbnailURI)
	if err != nil {
		log.Printf("Error inserting file into database: %v", err)
		return nil, err
	}

	log.Printf("File stored in database: ID = %d, URI = %s, ThumbnailURI = %s", file.ID, file.URI, file.ThumbnailURI)
	return &file, nil
}
