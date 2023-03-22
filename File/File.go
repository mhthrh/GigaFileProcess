package File

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

const (
	maxCapacity = 2 << 24
)

type File struct {
	file *os.File
}

func NewFile(path, name string) (*File, error) {
	file, err := os.Open(filepath.Join(path, name))
	if err != nil {
		return nil, fmt.Errorf("cannot open file %s: %w", name, err)
	}
	return &File{file: file}, nil
}

func (f *File) Read() ([]string, error) {

	var lines []string
	scanner := bufio.NewScanner(f.file)
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return []string{}, fmt.Errorf("cannot scan file %w", err)
	}
	return lines, nil
}

func (f *File) Close() error {
	return f.file.Close()
}
