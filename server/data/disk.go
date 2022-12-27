package data

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"server/logger"
	"strings"
)

func RemoveFromDisk(fileName string) {
	// logger.MyLog.Println(fileName)
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

// in case Go app goes down and redis service is still running and some data expired
func RemoveDeadRefs() {
	command := "ls | grep -E '([A-Za-z0-9]+(-[A-Za-z0-9]+)+)' -o | uniq"
	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = "uploaded/"
	out, err := cmd.Output()
	if err != nil {
		logger.MyLog.Println(err)
	}

	if len(out) != 0 {
		// fmt.Println(out)
		filesAfterTrimLastNewline := string(out)[:len(string(out))-1]
		files := strings.Split(filesAfterTrimLastNewline, "\n")
		// fmt.Println(files)
		for _, file := range files {
			if !InMemoryUUID.ItemExist(file) && file != "" {
				fmt.Println("dead refrences Found ", file)
				RemoveFromDisk(file)
			}
		}
	}
}
