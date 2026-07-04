package prompt

import (
	"fmt"
	"gh-cleaner/internal/bayes"
	"gh-cleaner/internal/files"
	gh "gh-cleaner/internal/github"
	"gh-cleaner/internal/structures"
	"log"
	"strings"

	"github.com/google/go-github/v62/github"
	"github.com/jbrukh/bayesian"
	"github.com/manifoldco/promptui"
)

// displayRepo prints repository details in a formatted manner
func displayRepo(repo *github.Repository) {
	fmt.Print("\n")
	fmt.Println("===================================================")
	fmt.Printf("repoName: %s\n", repo.GetName())
	fmt.Printf("repoDesc: %s\n", repo.GetDescription())
	fmt.Printf("repoLanguage: %s\n", repo.GetLanguage())
	if repo.GetLicense() != nil {
		fmt.Printf("repoLicense: %s\n", repo.GetLicense().GetKey())
	} else {
		fmt.Printf("repoLicense: %s\n", "None")
	}
	fmt.Printf("repoCreated: %s\n", repo.GetCreatedAt().Format("01-02-2006"))
	fmt.Printf("repoUpdated: %s\n", repo.GetUpdatedAt().Format("01-02-2006"))
	fmt.Printf("repoPushed: %s\n", repo.GetPushedAt().Format("01-02-2006"))
	fmt.Println("===================================================")
	fmt.Print("\n")
}

// convertToGitHubRepos converts a slice of structures.Repository to a slice of *github.Repository
func convertToGitHubRepos(srepos []*structures.Repository) []*github.Repository {
	var repos []*github.Repository
	for _, repo := range srepos {
		repos = append(repos, repo.OriginalRepo)
	}
	return repos
}

// promptDeleteRepo prompts the user to confirm the deletion of a repository
func promptDeleteRepo(repo *github.Repository) bool {
	displayRepo(repo)
	prompt := promptui.Prompt{
		Label: fmt.Sprintf("Can Delete the repo: %s [use (Y/N)]", repo.GetName()),
	}

	result, err := prompt.Run()
	if err != nil {
		return false
	}

	return strings.ToLower(result) == "y"
}

func SelectRepo(login structures.Login, dry_run bool, repos []*github.Repository, classifier *bayesian.Classifier, forkFlag bool) {
	if classifier == nil {
		log.Println("Classifier not initialized, skipping sort.")
		return
	}

	sortedRepos := bayes.SortRepos(repos, classifier)
	nrepos := convertToGitHubRepos(sortedRepos)

	for _, repo := range nrepos {
		if forkFlag || (!forkFlag && !repo.GetFork()) {
			if promptDeleteRepo(repo) {
				isConfirmed := confirmDeletion(repo)
				if err := files.SaveRepositoryFiles(repo, isConfirmed); err != nil {
					log.Println("Failed to save repository file: ", err.Error())
				}

				if !dry_run && isConfirmed {
					if err := gh.DeleteRepository(login, repo); err != nil {
						log.Println("Failed to delete repository: ", err.Error())
					}
				}
			} else {
				fmt.Println("Skipping...")
			}
		}
	}
}
