package files

import (
	"gh-cleaner/internal/structures"
	"os"
	"path/filepath"

	"github.com/google/go-github/v62/github"
)

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

func SaveRepositoryFiles(g *github.Repository, isDeleted bool) error {
	sub := "saved"
	if isDeleted {
		sub = "deleted"
	}

	dir := filepath.Join(repoDir(), sub)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return writeRepositoryFile(filepath.Join(dir, g.GetName()), g)
}
