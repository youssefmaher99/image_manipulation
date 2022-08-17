package data

import (
	"log"
	"os"
	"path/filepath"
)

func RemoveFromDisk(fileName string) {

	// find matching files in uploded directory
	uplodedFiles, err := findMatchingFiles("uploaded", fileName)
	handleErr(err)

	// remove matching files in uploded directory
	err = removeMatchingFiles(uplodedFiles)
	handleErr(err)

	// find matching files in filtered directory
	filteredFiles, err := findMatchingFiles("filtered", fileName)
	handleErr(err)

	// remove matching files in filtered directory
	err = removeMatchingFiles(filteredFiles)
	handleErr(err)

	// find matching files in archives directory
	archivesFiles, err := findMatchingFiles("archives", fileName)
	handleErr(err)

	// remove matching files in archives directory
	err = removeMatchingFiles(archivesFiles)
	handleErr(err)

}

func findMatchingFiles(dir string, fileName string) ([]string, error) {
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
		log.Println(err)
	}
}
