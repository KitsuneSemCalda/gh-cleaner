package main

import (
	"flag"
	"fmt"
	bayestheorem "gh-cleaner/internal/bayes_theorem"
	"gh-cleaner/internal/files"
	"gh-cleaner/internal/github"
	"gh-cleaner/internal/prompt"
	"gh-cleaner/internal/structures"
	"log"
)

var (
	netrcPath string = ""
	forks     bool   = true
	dry_run   bool   = false
	login     structures.Login
)

func main() {
	dry_run := flag.Bool("dry-run", false, "Enable dry run mode for training the bayes theorem")
	forks := flag.Bool("forks", false, "Enable forks to enable delete repository forked")
	flag.Parse()

	netrcPath = files.GetNetrc()

	if netrcPath == "" {
		log.Println(".netrc not founded")
		return
	}

	login = files.MountLogin(netrcPath)
	savedRepos, deletedRepos := files.GetInfoAboutRepo()

	if (savedRepos == nil) && (deletedRepos == nil) && (*dry_run == false) {
		fmt.Println("Run the code with flag --dry-run to create a mock files from this project")
		return
	}

	classifier := bayestheorem.GenerateClassifier(deletedRepos, savedRepos)

	prompt.SelectRepo(login, *dry_run, github.GetRepositoriesByToken(login), classifier, *forks)
}
