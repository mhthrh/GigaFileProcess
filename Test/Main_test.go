package Test

import (
	"os"
	"testing"
)

const (
	Count    = 100_000
	FileName = "testFile.csv"
)

var (
	FileArray []string
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
