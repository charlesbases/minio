package minio

import (
	"context"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	minioEndpoint       = "127.0.0.1:9000"
	minioEndpointSecure = false

	minioUser         = "minio"
	minioUserPassword = "minioadmin"

	contentType = "application/octet-stream"
)

var ctx = context.Background()

var client *minio.Client

func init() {
	NewClient()
}

// NewClient .
func NewClient() {
	cli, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioUser, minioUserPassword, ""),
		Secure: minioEndpointSecure,
	})
	if err != nil {
		log.Fatalln("minio.New() failed: ", err)
	}
	client = cli
}

// MakeBucket .
func MakeBucket(bucketName string, location string) {
	err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		if exists, _ := client.BucketExists(ctx, bucketName); exists {
			return
		}
		log.Fatalln("minio.MakeBucket() failed: ", err)
	}
}

// PutObject .
func PutObject(bucketName string, object io.Reader, objectName string, size int64) {
	_, err := client.PutObject(ctx, bucketName, objectName, object, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln("minio.PutObject() failed: ", err)
	}
}

// GetObject .
func GetObject(bucketName string, objectName string) io.Reader {
	obj, err := client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln("minio.GetObject() failed: ", err)
	}
	return obj
}
