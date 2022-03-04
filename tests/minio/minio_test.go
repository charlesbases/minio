package minio

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
)

const (
	bucketName     = "minio-test"
	bucketLocation = "us-east-1"
)

func init() {
	MakeBucket(bucketName, bucketLocation)
}

func TestMinioPutObject(t *testing.T) {
	// PutObjectName: 2022-03-01
	// PutObjectContent: 2022-03-01 09:00:00
	reader := strings.NewReader(current())
	PutObject(bucketName, reader, date(), reader.Size())
	log.Println("Successfully PutObject")
}

func TestMinioGetObject(t *testing.T) {
	// GetObjectName: 2022-03-01
	// SaveObjectPath: 2022-03-01 09:00:00
	obj := GetObject(bucketName, date())
	file, err := os.OpenFile(current(), os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatalln("os.OpenFile() failed: ", err)
	}
	defer file.Close()

	io.Copy(file, obj)
	log.Println("Successfully GetObject. dest: ", file.Name())
}

func TestMinioPutFolder(t *testing.T) {
	hook := func(folder string, amount int, size int) {
		for val := range serial(amount) {
			PutObject(bucketName, stream(size), filepath.Join(folder, val), int64(size))
		}
		log.Println("Successfully PutFolder ", folder)
	}

	var wg = sync.WaitGroup{}

	wg.Add(1)
	go func() {
		// PutObject 2022/03/01/{0000...9999} 1KiB
		hook("2022/03/01", 9999, 1<<10)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		// PutObject 2022/03/02/{0000...9999} 1KiB
		hook("2022/03/02", 9999, 1<<10)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		// PutObject 2022/03/03/{0000...9999} 1KiB
		hook("2022/03/03", 9999, 1<<10)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		// PutObject 2022/03/04/{0000...9999} 1MiB
		hook("2022/03/04", 9999, 1<<10)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		// PutObject 2022/03/05/{0000...9999} 1MiB
		hook("2022/03/05", 9999, 1<<10)
		wg.Done()
	}()

	wg.Wait()
}

func TestMinioGetFolder(t *testing.T) {
	// bucket list
	buckets, err := client.ListBuckets(ctx)
	if err != nil {
		log.Fatalln("minio.ListBuckets() failed: ", err)
	}

	for _, bucket := range buckets {
		var objects = make([]*minio.ObjectInfo, 0, 1024)
		for object := range client.ListObjects(ctx, bucket.Name, minio.ListObjectsOptions{
			Recursive: true,
		}) {
			objects = append(objects, &object)
		}
		log.Println(fmt.Sprintf("%s (%d Objects)", bucket.Name, len(objects)))
	}
}

// stream 生成指定大小的字节流
func stream(size int) io.Reader {
	var reader = new(bytes.Buffer)
	reader.Grow(size)
	for i := 0; i < size; i++ {
		reader.WriteByte(1)
	}
	return reader
}

func serial(maximum int) <-chan string {
	var ch = make(chan string, 1)
	go func() {
		var max = strconv.Itoa(maximum)
		for i := 0; i < maximum; i++ {
			var str = strconv.Itoa(i)
			ch <- strings.Repeat("0", len(max)-len(str)) + str
		}
		close(ch)
	}()
	return ch
}

func date() string {
	return time.Now().Format("2006-01-02")
}

func current() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
