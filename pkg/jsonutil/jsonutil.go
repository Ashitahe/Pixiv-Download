package jsonutil

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"pixivDownload/pkg/search"

	"github.com/tidwall/gjson"
)

func ReadJSON(filename string) (gjson.Result, error) {
    file, err := os.Open(filename)
	if err != nil {
        return gjson.Result{}, err
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        return gjson.Result{}, err
    }

    var result gjson.Result
    err = json.Unmarshal(data, &result)
    if err != nil {
        return gjson.Result{}, err
    }

    return result, nil
}

func ReadJSONFileToDownload(filename string) error {
	result, err := ReadJSON(filename)
	if err != nil {
		return err
	}

	result.ForEach(func(key, value gjson.Result) bool {
		err := search.SearchByIllustId(value.Str, "./downloads/")
		if err != nil {
			fmt.Printf("Failed to download image %s: %v\n", value.Str, err)
			return false // Stop iterating
		}
		return true // Continue iterating
	})

	return nil
}