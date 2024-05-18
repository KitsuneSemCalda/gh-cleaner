package files

import (
	"bufio"
	"gh-cleaner/internal/structures"
	"log"
	"os"
	"strings"
)

// This function receives the .netrc path and build a struct Login
func MountLogin(p string) structures.Login {
	var login, token string
	file, err := os.Open(p)
	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "machine github.com") {
			words := strings.Split(line, " ")
			login = words[3]
			token = words[5]
		}
	}

	return structures.CreateLogin(login, token)
}
