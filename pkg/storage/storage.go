package storage

import (
	"bufio"
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

	writer := bufio.NewWriter(file)
	reader := bytes.NewReader(content)
	buffer := make([]byte, 32*1024)

	for {
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := writer.Write(buffer[:n]); err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	fmt.Printf("%s downloaded\n", filename)
	return nil
}
