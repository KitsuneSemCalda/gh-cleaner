package files

import (
	"bufio"
	"gh-cleaner/internal/structures"
	"log"
	"os"
	"path/filepath"
)

func loadRepositories(dir string) ([]structures.RepoDatum, error) {
	var repos []structures.RepoDatum

	spath := filepath.Join(repoDir(), dir)

	entries, err := os.ReadDir(spath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		filePath := filepath.Join(spath, entry.Name())
		f, err := os.Open(filePath)
		if err != nil {
			log.Println("Occurred an unknown error in open file: ", err.Error())
			continue
		}

		var lines []string
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		f.Close()

		if err := scanner.Err(); err != nil {
			log.Println("Occurred an unknown error in scan file: ", err.Error())
			continue
		}

		if len(lines) >= 5 {
			repos = append(repos, structures.RepoDatum{
				Name:        lines[0],
				Stars:       lines[1],
				Issues:      lines[2],
				Forks:       lines[3],
				Subscribers: lines[4],
			})
		}
	}

	return repos, nil
}

func GetInfoAboutRepo() ([]structures.RepoDatum, []structures.RepoDatum, error) {
	saved, err := loadRepositories("saved")
	if err != nil {
		return nil, nil, err
	}

	deleted, err := loadRepositories("deleted")
	if err != nil {
		return nil, nil, err
	}

	return saved, deleted, nil
}
