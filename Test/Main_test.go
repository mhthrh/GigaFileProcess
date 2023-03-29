package Test

import (
	"database/sql"
	_ "github.com/godror/godror"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const (
	ftpIp                 = "localhost"
	ftpPort               = 21
	ftpUser               = "ftpuser"
	ftpPassword           = "123456"
	Count                 = 100_000
	FileName              = "testFile.csv"
	fileTestDirector      = "Files"
	filetestName4Upload   = "foo.zip"
	filetestName4Download = "foo1.zip"
)

var (
	FileArray []string
	cnn       *sql.DB
	err       error
)

func init() {
	cnn, err = sql.Open("godror", `user="mohsen" password="mohsen" connectString="localhost:1521/xe"
    poolSessionTimeout=42s configDir="/tmp/admin"
    heterogeneousPool=false standaloneConnection=false`)
	if err != nil {
		log.Fatalln("cannot connect to db")
	}
	f, err := os.Create(filepath.Join(fileTestDirector, filetestName4Upload))
	if err != nil {
		log.Fatal(err)
	}

	if err := f.Truncate(1e7); err != nil {
		log.Fatal(err)
	}
}
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
