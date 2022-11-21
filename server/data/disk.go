package data

import (
	"os"
	"path/filepath"
	"server/logger"
)

func RemoveFromDisk(fileName string) {
	logger.MyLog.Println(fileName)
	// find matching files in uploded directory
	uplodedFiles, err := findMatchingFilesInDirs("uploaded", fileName)
	handleErr(err)

	// remove matching files in uploded directory
	err = removeMatchingFiles(uplodedFiles)
	handleErr(err)

	// find matching files in filtered directory
	filteredFiles, err := findMatchingFilesInDirs("filtered", fileName)
	handleErr(err)

	// remove matching files in filtered directory
	err = removeMatchingFiles(filteredFiles)
	handleErr(err)

	// find matching files in archives directory
	archivesFiles, err := findMatchingFilesInDirs("archives", fileName)
	handleErr(err)

	// remove matching files in archives directory
	err = removeMatchingFiles(archivesFiles)
	handleErr(err)

}

func findMatchingFilesInDirs(dir string, fileName string) ([]string, error) {
	files, err := filepath.Glob(dir + "/" + fileName + "*")
	if err != nil {
		return nil, err
	}
	return files, nil
}

func removeMatchingFiles(files []string) error {
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}

func handleErr(err error) {
	if err != nil {
		logger.MyLog.Println(err)
	}
}
