package files

import (
	"os"
	"path/filepath"
)

// This function returns the user home
// we use the os.GetEnv to get the enviroment variable HOME
func getHome() string {
	return os.Getenv("HOME")
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
