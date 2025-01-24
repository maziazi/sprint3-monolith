package awsr

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"os"
	"path/filepath"
)

func UploadToS3(filePath string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-2"),
	})
	if err != nil {
		log.Printf("Error creating AWS session: %v", err)
		return "", fmt.Errorf("unable to create AWS session: %v", err)
	}

	s3Client := s3.New(sess)

	log.Println("Opening file:", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file %v: %v", filePath, err)
		return "", fmt.Errorf("unable to open file %v: %v", filePath, err)
	}
	defer file.Close()
	log.Println("File opened successfully.")

	fileName := filepath.Base(filePath)

	log.Printf("Uploading file %v to S3...", fileName)

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String("penopangsistemuii"), //TODO EDIT INI
		Key:         aws.String(fileName),
		Body:        file,
		ContentType: aws.String("application/octet-stream"),
		ACL:         aws.String("public-read"),
	})
	if err != nil {
		log.Printf("Error uploading %v to S3: %v", fileName, err)
		return "", err
	}
	log.Printf("File %v uploaded to S3 successfully.", fileName)
	//TODO EDIT INI JUGA BAGIAN FILE URL
	fileURL := fmt.Sprintf("https://penopangsistemuii.s3.ap-southeast-2.amazonaws.com/%s", fileName)
	log.Printf("File URL: %v", fileURL)

	return fileURL, nil
}
