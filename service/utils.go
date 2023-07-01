package service

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

func gzippedBytes(data []byte) (gzipped []byte, err error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	_, err = zw.Write(data)
	if err != nil {
		return []byte{}, err
	}
	err = zw.Close()
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func fileSha1(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	return ioReaderSha1(file)
}

func bytesSha1(data []byte) string {
	reader := bytes.NewReader(data)
	return ioReaderSha1(reader)
}

func ioReaderSha1(reader io.Reader) string {
	hash := sha1.New()
	if _, err := io.Copy(hash, reader); err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}
