package storage

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func SaveFile(content []byte, directory, filename string) error {
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		fmt.Printf("Error: can't create directory, using current directory instead: %v\n", err)
		directory = "./"
	}

	fullPath := filepath.Join(directory, filename)
	file, err := os.Create(fullPath)
	if err != nil {
		fmt.Printf("Failed to create file %s: %v\n", fullPath, err)
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, bytes.NewReader(content)); err != nil {
		fmt.Printf("Failed to write to file %s: %v\n", fullPath, err)
		return err
	}

	fmt.Printf("%s downloaded\n", filename)
	return nil
}
