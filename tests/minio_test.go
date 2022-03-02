package tests

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	endpoint = "127.0.0.1:9000"
	secure   = false

	user     = "minio"
	password = "minioadmin"
)

var (
	client *minio.Client
	ctx    = context.Background()
)

func TestMinioClient(t *testing.T) {
	cli, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(user, password, ""),
		Secure: secure,
	})
	if err != nil {
		log.Fatalln("minio.New() failed: ", err)
	}

	client = cli
}

const (
	bucketName = "minio-test"
	location   = "us-east-1"
)

func TestMinioBucket(t *testing.T) {
	TestMinioClient(t)

	err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		if exists, _ := client.BucketExists(ctx, bucketName); !exists {
			log.Fatalln("minio.MakeBucket() failed: ", err)
		}
	} else {
		log.Fatalln("minio.MakeBucket() failed: ", err)
	}
}

const (
	objectName  = "2022-03-01"
	objectPath  = "2022-03-01"
	contentType = "application/octet-stream"
)

func TestMinioUpload(t *testing.T) {
	TestMinioBucket(t)

	info, err := client.FPutObject(ctx, bucketName, objectName, objectPath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln("minio.FPutObject() failed: ", err)
	}
	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
}

var destPath = "2022-03-01 " + time.Now().Format("15:04:05")

func TestMinioDownload(t *testing.T) {
	TestMinioBucket(t)

	err := client.FGetObject(ctx, bucketName, objectName, destPath, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln("minio.FGetObject() failed: ", err)
	}
}
