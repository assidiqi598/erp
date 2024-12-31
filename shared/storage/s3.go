package storage

import (
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3ClientType struct {
	*minio.Client
}

var S3Client *S3ClientType

type S3Credentials struct {
	Endpoint  string
	AccessKey string
	SecretKey string
}

func CreateS3Client(cred S3Credentials) error {
	client, err := minio.New(cred.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cred.AccessKey, cred.SecretKey, ""),
		Secure: true,
	})

	if err != nil {
		return fmt.Errorf("failed to initialize S3 client: %v", err)
	}

	log.Println("Connected to S3")

	S3Client.Client = client

	return nil
}

func GetS3Client() *S3ClientType {
	return S3Client
}
