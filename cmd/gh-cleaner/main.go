package main

import (
	"flag"
	"fmt"
	"gh-cleaner/internal/bayes"
	"gh-cleaner/internal/files"
	"gh-cleaner/internal/github"
	"gh-cleaner/internal/prompt"
	"gh-cleaner/internal/structures"
	"log"

	"github.com/jbrukh/bayesian"
)

func main() {
	dryRun := flag.Bool("dry-run", false, "Enable dry run mode for training the bayes theorem")
	forks := flag.Bool("forks", false, "Enable forks to enable delete repository forked")
	flag.Parse()

	login, err := loadLogin()
	if err != nil {
		log.Println(err.Error())
		return
	}

	classifier := buildClassifier()

	repos, err := github.GetRepositoriesByToken(login)
	if err != nil {
		log.Println("Error fetching repositories: ", err.Error())
		return
	}

	prompt.SelectRepo(login, *dryRun, repos, classifier, *forks)
}

func loadLogin() (structures.Login, error) {
	netrcPath := files.GetNetrc()
	if netrcPath == "" {
		return structures.Login{}, fmt.Errorf(".netrc not found")
	}

	login, err := files.MountLogin(netrcPath)
	if err != nil {
		return structures.Login{}, fmt.Errorf("error reading .netrc: %w", err)
	}

	return login, nil
}

func buildClassifier() *bayesian.Classifier {
	savedRepos, deletedRepos, err := files.GetInfoAboutRepo()
	if err != nil {
		log.Printf("Warning: could not load repository data: %v", err)
	}
	return bayes.GenerateClassifier(deletedRepos, savedRepos)
}
