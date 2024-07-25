package files

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
)

func loadRepositorySaved() [][]string {
	var RepositorySaved [][]string

	spath := filepath.Join(getHome(), ".local", "share", "gh-cleaner", "repository", "saved")

	files, err := os.ReadDir(spath)
	if err != nil {
		log.Println("Occured an unknown error in read dir: ", err.Error())
	}

	for _, file := range files {
		filePath := filepath.Join(spath, file.Name())
		f, err := os.Open(filePath)
		if err != nil {
			log.Println("Occured an unknown error in open file: ", err.Error())
		}
		defer f.Close()

		var lines []string
		scanner := bufio.NewScanner(f)
		// Read file line by line
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			log.Println("Occured an unknown error in open file: ", err.Error())
			return nil
		}

		RepositorySaved = append(RepositorySaved, lines)
	}

	return RepositorySaved
}

func loadRepositoryDeleted() [][]string {
	var RepositoryDeleted [][]string

	spath := filepath.Join(getHome(), ".local", "share", "gh-cleaner", "repository", "deleted")

	files, err := os.ReadDir(spath)
	if err != nil {
		log.Println("Occured an unknown error in read dir: ", err.Error())
	}

	for _, file := range files {
		filePath := filepath.Join(spath, file.Name())
		f, err := os.Open(filePath)
		if err != nil {
			log.Println("Occured an unknown error in open file: ", err.Error())
		}
		defer f.Close()

		var lines []string
		scanner := bufio.NewScanner(f)
		// Read file line by line
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			log.Println("Occured an unknown error in open file: ", err.Error())
			return nil
		}

		RepositoryDeleted = append(RepositoryDeleted, lines)
	}

	return RepositoryDeleted
}

func GetInfoAboutRepo() ([][]string, [][]string) {
	return loadRepositorySaved(), loadRepositoryDeleted()
}
