package Test

import (
	"fmt"
	"github.com/mhthrh/GigaFileProcess/Utils/Random"
	File "github.com/mhthrh/GigaFileProcess/file"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestFile(t *testing.T) {
	writeFile(t)
	readFile(t)
}
func writeFile(t *testing.T) {
	file, err := os.Create(FileName)
	require.NoError(t, err)
	require.NotEmpty(t, file)
	var sb strings.Builder

	for i := 0; i < Count; i++ {
		sb.WriteString(fmt.Sprintf("%d,%s,%s,%s,%s,%d\n", i, Random.RandomString(10), Random.RandomString(10), Random.RandomString(10), Random.RandomString(10), Random.RandomInt(2, 10)))
	}

	_, err = file.WriteString(sb.String())
	require.NoError(t, err)

}
func readFile(t *testing.T) {
	csv, err := File.NewFile("./", FileName)
	require.NoError(t, err)
	require.NotEmpty(t, csv)

	FileArray, err = csv.Read()
	require.NoError(t, err)
	require.Equal(t, len(FileArray), Count)
	err = csv.Close()
	require.NoError(t, err)
}

//func TestDoProcess(t *testing.T) {
//	err := FileProcess.DoProcess(FileArray)
//	require.NoError(t, err)
//}
