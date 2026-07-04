package files

import (
	"gh-cleaner/internal/structures"
	"os"
	"path/filepath"

	"github.com/google/go-github/v62/github"
)

func createDirectory(path string) {
	if !fileExists(path) {
		os.MkdirAll(path, 0755)
	}
}

func writeRepositoryFile(fpath string, g *github.Repository) error {
	structure := structures.CreateRepository(g)
	file, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, element := range structure.DataFields() {
		if _, err := file.WriteString(element + "\n"); err != nil {
			return err
		}
	}

	return nil
}

func createFile(path string, file string, g *github.Repository) error {
	fpath := filepath.Join(path, file)
	return writeRepositoryFile(fpath, g)
}

func SaveRepositoryFiles(g *github.Repository, isDeleted bool) error {
	path := repoDir()

	if isDeleted {
		deletedPath := filepath.Join(path, "deleted")
		createDirectory(deletedPath)
		return createFile(deletedPath, g.GetName(), g)
	}

	savedPath := filepath.Join(path, "saved")
	createDirectory(savedPath)
	return createFile(savedPath, g.GetName(), g)
}
