package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/mhthrh/GigaFileProcess/Utils/CryptoUtil"
	"github.com/mhthrh/GigaFileProcess/Utils/Random"
	"github.com/mhthrh/GigaFileProcess/entity"
	"github.com/mhthrh/GigaFileProcess/ftp"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	ftpIp       = "localhost"
	ftpPort     = 21
	ftpUser     = "ftpuser"
	ftpPassword = "123456"
)

func main() {

	postBody, err := json.Marshal(request())
	if err != nil {
		log.Fatalln(err)
	}
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:8569/gateway/file", "application/json", responseBody)

	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(resp)
}

func request() entity.FileRequest {

	const amount = 3
	path := "client/files"
	filename := Random.RandomString(10)
	defer os.Remove(filename)
	count := int(Random.RandomInt(1000, 1000000))
	file, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		log.Fatalln(err)
	}
	var sb strings.Builder
	for i := 0; i < count; i++ {
		sb.WriteString(fmt.Sprintf("%d,%s,%s,%s,%s,%d\n", i, Random.RandomString(10), Random.RandomString(10), Random.RandomString(10), Random.RandomString(10), amount))
	}
	_, err = file.WriteString(sb.String())
	if err != nil {
		log.Fatalln(err)
	}
	k := CryptoUtil.NewKey()
	k.FilePath = filepath.Join(path, filename)
	md5, _ := k.Md5Sum()

	client, err := ftp.New(ftpIp, ftpUser, ftpPassword, ftpPort)
	if err != nil {
		log.Fatalln(err)
	}
	err = client.Upload(filename, filename, path)
	if err != nil {
		log.Fatalln(err)
	}
	return entity.FileRequest{
		ID:       uuid.New(),
		Md5:      md5,
		Name:     filename,
		Count:    count,
		Sum:      float64(count * 3),
		Priority: 30,
	}
}
