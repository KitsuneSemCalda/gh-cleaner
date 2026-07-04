package files

import (
	"os"
	"path/filepath"
)

func getHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return os.Getenv("HOME")
	}
	return home
}

// This function check if the file fileExists
// We pass the path of file and check if stats is unequal of IsNotExist
func fileExists(f string) bool {
	_, err := os.Stat(f)

	return !os.IsNotExist(err)
}

// This function generates the default path from .netrc and return if file exists
func GetNetrc() string {
	path := filepath.Join(getHome(), ".netrc")

	if fileExists(path) {
		return path
	}

	return ""
}
