package munichain

import (
	"encoding/json"
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
	blocks := getBlocksFilePath(dataDir)
	_, err := os.Stat(blocks)
	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}

	dbDir := getDbDirPath(dataDir)
	if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
		return err
	}
	if err := writeGenesisBlockToBlocksFile(blocks); err != nil {
		return err
	}
	return nil
}

func writeGenesisBlockToBlocksFile(path string) error {
	hash, err := GenesisBlock.Hash()
	if err != nil {
		return err
	}
	genesisJson, err := json.Marshal(BlockFS{Key: hash, Value: GenesisBlock})
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, append(genesisJson, '\n'), os.ModePerm)
}
