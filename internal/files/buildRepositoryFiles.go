package files

import (
	"gh-cleaner/internal/structures"
	"log"
	"os"
	"path/filepath"

	"github.com/google/go-github/v62/github"
)

func createDirectory(path string) {
	if !fileExists(path) {
		os.MkdirAll(path, 0777)
	}
}

func populateUnexistFile(fpath string, g *github.Repository) {
	structure := structures.CreateRepository(g)
	file, err := os.Create(fpath)
	if err != nil {
		log.Println("Ocurred an unknown error in create file: ", err.Error())
		return
	}

	for _, element := range structure.GetAllValues() {
		file.WriteString(element)
	}
}

func populateExistedFile(fpath string, g *github.Repository) {
	structure := structures.CreateRepository(g)
	file, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		log.Println("Ocurred an unknown error in create file: ", err.Error())
		return
	}

	for _, element := range structure.GetAllValues() {
		file.WriteString(element + "\n")
	}
}

func createFile(path string, file string, g *github.Repository) {
	fpath := filepath.Join(path, file)
	if !fileExists(fpath) {
		populateUnexistFile(fpath, g)
	}

	populateExistedFile(fpath, g)
}

func SaveRepositoryFiles(g *github.Repository, isDeleted bool) {
	path := filepath.Join(getHome(), ".local", "share", "gh-cleaner", "repository")

	if isDeleted {
		deletedPath := filepath.Join(path, "deleted")

		createDirectory(deletedPath)
		createFile(deletedPath, g.GetName(), g)
		return
	}

	savedPath := filepath.Join(path, "saved")
	createDirectory(savedPath)
	createFile(savedPath, g.GetName(), g)
}
