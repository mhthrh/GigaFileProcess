package Test

import (
	"database/sql"
	_ "github.com/godror/godror"
	"log"
	"os"
	"testing"
)

const (
	Count    = 100_000
	FileName = "testFile.csv"
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
}
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
