package munichain

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

func loadJson(v interface{}, pathChunks ...string) error {
	file, err := loadFile(pathChunks...)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, v)
}

func loadFile(pathChunks ...string) ([]byte, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	pathChunks = append([]string{cwd}, pathChunks...)
	path := filepath.Join(pathChunks...)

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func writeJson(v interface{}, pathChunks ...string) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return writeFile(data, pathChunks...)
}

func writeFile(data []byte, pathChunks ...string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	pathChunks = append([]string{cwd}, pathChunks...)
	path := filepath.Join(pathChunks...)
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
}
