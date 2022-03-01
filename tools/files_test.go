package tools

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"
)

var files = map[string]int{
	"1-1KiB":  1 << 10,
	"2-2KiB":  1 << 11,
	"3-4KiB":  1 << 12,
	"4-8KiB":  1 << 13,
	"5-16KiB": 1 << 14,
	"6-1MiB":  1 << 20,
	"7-4MiB":  1 << 22,
	"8-1GiB":  1 << 30,
	"9-4GiB":  1 << 32,
}

func Test(t *testing.T) {
	for name, size := range files {
		create(name, size)
	}
}

func create(name string, size int) {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	io.Copy(file, useless(size))
}

func useless(size int) io.Reader {
	var reader = new(bytes.Buffer)
	reader.Grow(size)
	for i := 0; i < size; i++ {
		reader.WriteByte(1)
	}
	return reader
}
