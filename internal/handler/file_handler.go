package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"sprint3/internal/awsr"
	"sprint3/internal/service"
	"strconv"
	"strings"
)

func UploadFileHandler(c *gin.Context) {
	// Mengambil file dari form-data
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not found"})
		return
	}

	// Memvalidasi ekstensi file
	fileExtension := strings.ToLower(filepath.Ext(file.Filename))
	if fileExtension != ".jpeg" && fileExtension != ".jpg" && fileExtension != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only jpeg, jpg, png allowed."})
		return
	}

	// Memvalidasi ukuran file (maksimum 100KB)
	if file.Size > 1024*100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 100KiB"})
		return
	}

	// Menyimpan file sementara secara lokal
	uploadPath := filepath.Join("uploads", file.Filename)
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file locally"})
		return
	}

	// Upload ke S3
	fileURL, err := awsr.UploadToS3(uploadPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to S3"})
		return
	}

	// Membuat thumbnail file
	thumbnailURL, err := awsr.CreateThumbnailAndUploadToS3(uploadPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create thumbnail"})
		return
	}

	// Menyimpan data file ke database
	storedFile, err := service.AddFile(fileURL, thumbnailURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store file in the database"})
		return
	}

	// Menyusun response
	c.JSON(http.StatusOK, gin.H{
		"fileId":           strconv.Itoa(storedFile.ID),
		"fileUri":          storedFile.URI,
		"fileThumbnailUri": storedFile.ThumbnailURI,
	})
}
