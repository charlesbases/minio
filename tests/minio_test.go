package tests

import (
	"context"
	"io"
	"log"
	"os"
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
	minioClient()
}

func TestMinioMakeBucket(t *testing.T) {
	minioMakeBucket()
}

func TestMinioPutObject(t *testing.T) {
	minioMakeBucket()

	// 获取字节流
	file, err := os.Open(objectPath)
	if err != nil {
		log.Fatalln("os.Open() failed: ", err)
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		log.Fatalln("file.Stat() failed: ", err)
	}

	minioPutObject(file, fileStat.Size())
}

var destPath = "2022-03-01 " + time.Now().Format("15:04:05")

func TestMinioGetObject(t *testing.T) {
	minioMakeBucket()

	obj := minioGetObject()
	file, err := os.OpenFile(destPath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatalln("os.OpenFile() failed: ", err)
	}
	defer file.Close()

	io.Copy(file, obj)
}

func minioClient() {
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

func minioMakeBucket() {
	minioClient()

	err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		if exists, _ := client.BucketExists(ctx, bucketName); exists {
			return
		}
		log.Fatalln("minio.MakeBucket() failed: ", err)
	}
}

const (
	objectName  = "2022-03-01"
	objectPath  = "2022-03-01"
	contentType = "application/octet-stream"
)

func minioPutObject(obj io.Reader, size int64) {
	info, err := client.PutObject(ctx, bucketName, objectName, obj, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln("minio.PutObject() failed: ", err)
	}
	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
}

func minioGetObject() *minio.Object {
	obj, err := client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln("minio.FGetObject() failed: ", err)
	}
	return obj
}
