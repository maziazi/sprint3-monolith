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
func CreateThumbnailAndUploadToS3(filePath string) (string, error) {
	// Membaca file image
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Meng-decode gambar
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// Resize gambar menjadi thumbnail (misalnya 100x100px)
	thumb := resize.Thumbnail(100, 100, img, resize.Lanczos3)

	// Menyimpan thumbnail ke buffer
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, thumb, nil)
	if err != nil {
		return "", err
	}

	// Menyimpan thumbnail ke S3
	thumbnailFileName := GenerateUniqueFileName("jpg")
	thumbnailURL, err := uploadToS3(buf.Bytes(), thumbnailFileName)
	if err != nil {
		return "", err
	}

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
