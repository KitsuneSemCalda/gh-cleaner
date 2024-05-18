package main

import (
	"gh-cleaner/internal/files"
	"gh-cleaner/internal/github"
	"gh-cleaner/internal/prompt"
	"gh-cleaner/internal/structures"
)

var (
	netrcPath string = ""
	forks     bool   = true
	login     structures.Login
)

func main() {
	netrcPath = files.GetNetrc()

	if netrcPath == "" {
		return
	}

	login = files.MountLogin(netrcPath)
	prompt.SelectRepo(login, github.GetRepositoriesByToken(login), forks)
}
