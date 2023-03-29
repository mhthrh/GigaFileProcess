package ftp

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/secsy/goftp"
	"io"
	"os"
	"path/filepath"
	"time"
)

type FtpClient struct {
	client *goftp.Client
}

func New(ip, username, password string, port int) (*FtpClient, error) {
	client, err := goftp.DialConfig(goftp.Config{
		User:     username,
		Password: password,
		Timeout:  time.Second * 5,
	}, fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return nil, fmt.Errorf("cannot connect to ftp server, %v", err)
	}
	return &FtpClient{
		client: client,
	}, nil
}
func (f *FtpClient) CreateDirectory(path string) error {
	_, err := f.client.Mkdir(path)
	return err
}
func (f *FtpClient) DeleteDirectory(path string) error {
	return f.client.Rmdir(path)
}
func (f *FtpClient) List(path string) ([]string, error) {
	var result []string
	files, err := f.client.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("an error occured on reading directory, %v", err)
	}
	for _, file := range files {
		result = append(result, file.Name())
	}
	return result, nil
}
func (f *FtpClient) Download(downloadFileName, localFileName, pathDestination string) error {
	file, err := os.Create(filepath.Join(pathDestination, localFileName))
	if err != nil {
		return fmt.Errorf("cannot create file, %v", err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	err = f.client.Retrieve(downloadFileName, writer)

	if err != nil {
		return fmt.Errorf("cannot write file, %v", err)
	}
	return writer.Flush()

}
func (f *FtpClient) Upload(filename, ftpPath, path string) error {
	file, err := os.Open(filepath.Join(path, filename))
	if err != nil {
		return fmt.Errorf("cannot find source file")
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	part := make([]byte, 1024)
	var count int
	for {
		if count, err = reader.Read(part); err != nil {
			break
		}
		buffer.Write(part[:count])
	}
	if err != io.EOF {
		return fmt.Errorf("buffering error, %v", err)
	}

	return f.client.Store(ftpPath, buffer)
}
func (f *FtpClient) Close() error {
	return f.client.Close()
}
