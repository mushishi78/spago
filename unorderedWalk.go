package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func unorderedWalk(base string, relativePath string, walkFn func(fileName string)) error {
	dirFile, err := os.Open(filepath.Join(base, relativePath))
	if err != nil {
		return fmt.Errorf("failed to open dir for walking: %v", err)
	}
	defer dirFile.Close()

	fileInfos, err := dirFile.Readdir(-1)
	if err != nil {
		return fmt.Errorf("failed to get info for files in dir: %v", err)
	}

	for _, fileInfo := range fileInfos {
		fileRelativePath := path.Join(relativePath, fileInfo.Name())

		if !fileInfo.IsDir() {
			walkFn(fileRelativePath)
			continue
		}

		err := unorderedWalk(base, fileRelativePath, walkFn)
		if err != nil {
			return err
		}
	}

	return nil
}
