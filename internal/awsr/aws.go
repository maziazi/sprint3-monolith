package awsr

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"time"
)

// GenerateUniqueFileName menghasilkan nama file unik menggunakan timestamp dan UUID
func GenerateUniqueFileName(ext string) string {
	// Menggunakan timestamp dan UUID untuk memastikan nama file unik
	return fmt.Sprintf("%d_%s.%s", time.Now().Unix(), uuid.New().String(), ext)
}

// UploadToS3 mengunggah file ke S3 dan mengembalikan URL file yang diunggah.
func UploadToS3(filePath string) (string, error) {
	// Membuka session AWS
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-2"), // Sesuaikan dengan region Anda
	})
	if err != nil {
		log.Printf("Error creating AWS session: %v", err)
		return "", fmt.Errorf("unable to create AWS session: %v", err)
	}

	// Membuat client S3
	s3Client := s3.New(sess)

	// Membuka file
	log.Println("Opening file:", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file %v: %v", filePath, err)
		return "", fmt.Errorf("unable to open file %v: %v", filePath, err)
	}
	defer file.Close()
	log.Println("File opened successfully.")

	// Mengambil ekstensi file
	fileExtension := filepath.Ext(filePath)
	// Menghasilkan nama file unik
	uniqueFileName := GenerateUniqueFileName(fileExtension)

	// Mengunggah file ke S3
	log.Printf("Uploading file %v to S3...", uniqueFileName)
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String("penopangsistemuii"), // Ganti dengan nama bucket Anda
		Key:         aws.String(uniqueFileName),
		Body:        file,
		ContentType: aws.String("application/octet-stream"),
		ACL:         aws.String("public-read"),
	})
	if err != nil {
		log.Printf("Error uploading %v to S3: %v", uniqueFileName, err)
		return "", err
	}
	log.Printf("File %v uploaded to S3 successfully.")

	// Mengembalikan URL file setelah diupload
	fileURL := fmt.Sprintf("https://penopangsistemuii.s3.ap-southeast-2.amazonaws.com/%s", uniqueFileName)
	log.Printf("File URL: %v", fileURL)

	return fileURL, nil
}

// CreateThumbnailAndUploadToS3 membuat thumbnail dari file gambar dan mengunggahnya ke S3.
// CreateThumbnailAndUploadToS3 membuat thumbnail dari file gambar dan mengunggahnya ke S3.
func CreateThumbnailAndUploadToS3(filePath string) (string, error) {
	log.Printf("Starting thumbnail creation for file: %s", filePath)

	// Membuka file gambar
	log.Println("Opening the image file...")
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file %s: %v", filePath, err)
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()
	log.Println("File opened successfully.")

	// Decode gambar
	log.Println("Decoding the image...")
	img, format, err := image.Decode(file)
	if err != nil {
		log.Printf("Error decoding the image: %v", err)
		return "", fmt.Errorf("failed to decode image: %v", err)
	}
	log.Printf("Image decoded successfully. Format: %s", format)

	// Membuat thumbnail (ukuran 100x100px)
	log.Println("Resizing the image to create a thumbnail...")
	thumb := resize.Thumbnail(100, 100, img, resize.Lanczos3)
	if thumb == nil {
		log.Println("Failed to resize the image.")
		return "", fmt.Errorf("failed to resize the image")
	}
	log.Println("Thumbnail created successfully.")

	// Menyimpan thumbnail ke buffer
	log.Println("Encoding the thumbnail to JPEG format...")
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, thumb, nil)
	if err != nil {
		log.Printf("Error encoding thumbnail to JPEG: %v", err)
		return "", fmt.Errorf("failed to encode thumbnail: %v", err)
	}
	log.Printf("Thumbnail encoded successfully. Buffer size: %d bytes", buf.Len())

	// Menghasilkan nama file unik untuk thumbnail
	thumbnailFileName := GenerateUniqueFileName("jpg")
	log.Printf("Generated unique filename for thumbnail: %s", thumbnailFileName)

	// Mengunggah thumbnail ke S3
	log.Printf("Uploading thumbnail %s to S3...", thumbnailFileName)
	thumbnailURL, err := uploadToS3(buf.Bytes(), thumbnailFileName)
	if err != nil {
		log.Printf("Error uploading thumbnail %s to S3: %v", thumbnailFileName, err)
		return "", fmt.Errorf("failed to upload thumbnail: %v", err)
	}
	log.Printf("Thumbnail uploaded successfully. URL: %s", thumbnailURL)

	return thumbnailURL, nil
}

// uploadToS3 mengunggah data byte ke S3 dan mengembalikan URL file yang diunggah.
func uploadToS3(fileData []byte, fileName string) (string, error) {
	// Membuka session AWS
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-2"), // Sesuaikan dengan region Anda
	})
	if err != nil {
		log.Printf("Error creating AWS session: %v", err)
		return "", fmt.Errorf("unable to create AWS session: %v", err)
	}

	// Membuat client S3
	s3Client := s3.New(sess)

	// Mengunggah file ke S3
	log.Printf("Uploading file %v to S3...", fileName)
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String("penopangsistemuii"), // Ganti dengan nama bucket Anda
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(fileData),
		ContentType: aws.String("application/octet-stream"),
		ACL:         aws.String("public-read"),
	})
	if err != nil {
		log.Printf("Error uploading %v to S3: %v", fileName, err)
		return "", err
	}
	log.Printf("File %v uploaded to S3 successfully.")

	// Mengembalikan URL file setelah diupload
	fileURL := fmt.Sprintf("https://penopangsistemuii.s3.ap-southeast-2.amazonaws.com/%s", fileName)
	log.Printf("File URL: %v", fileURL)

	return fileURL, nil
}
