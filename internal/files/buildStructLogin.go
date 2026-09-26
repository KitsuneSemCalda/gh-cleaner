package files

import (
	"bufio"
	"fmt"
	"gh-cleaner/internal/structures"
	"os"
	"strings"
)

func MountLogin(p string) (structures.Login, error) {
	var login, token string
	found := false
	file, err := os.Open(p)
	if err != nil {
		return structures.Login{}, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "machine github.com") {
			found = true
			words := strings.Fields(line)
			for i, w := range words {
				if w == "login" && i+1 < len(words) {
					login = words[i+1]
				}
				if w == "password" && i+1 < len(words) {
					token = words[i+1]
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return structures.Login{}, fmt.Errorf("error reading %s: %w", p, err)
	}

	if !found {
		return structures.Login{}, fmt.Errorf("no \"machine github.com\" entry found in %s", p)
	}

	return structures.CreateLogin(login, token), nil
}
