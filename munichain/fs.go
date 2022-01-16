package munichain

import (
	"io/ioutil"
	"os"
	"path"
)

func getDbDirPath(dataDir string) string {
	return path.Join(dataDir, "db")
}

func getBlocksFilePath(dataDir string) string {
	return path.Join(getDbDirPath(dataDir), "blocks.db")
}

func initDataDirIfNotExists(dataDir string) error {
	dbDir := getDbDirPath(dataDir)
	if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
		return err
	}
	blocks := getBlocksFilePath(dataDir)
	if err := writeEmptyBlocksFileToDisk(blocks); err != nil {
		return err
	}
	return nil
}

func writeEmptyBlocksFileToDisk(path string) error {
	return ioutil.WriteFile(path, []byte(""), os.ModePerm)
}
