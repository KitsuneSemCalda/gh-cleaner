package files

import "path/filepath"

func repoDir() string {
	return filepath.Join(getHome(), ".local", "share", "gh-cleaner", "repository")
}
