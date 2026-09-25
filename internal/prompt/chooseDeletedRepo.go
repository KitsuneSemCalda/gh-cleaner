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

// promptDeleteRepoResult interprets the free-text Y/N prompt. Separated from
// promptDeleteRepo so a cancelled/failed prompt (err != nil) can be told
// apart from a real "keep this repo" answer: only the latter is a decision
// worth teaching the classifier.
func promptDeleteRepoResult(result string, err error) (wantsDelete bool, cancelled bool) {
	if err != nil {
		return false, true
	}
	return strings.ToLower(result) == "y", false
}

// promptDeleteRepo prompts the user to confirm the deletion of a repository.
// cancelled is true when the prompt itself failed or was interrupted, as
// opposed to the user deliberately answering "no".
func promptDeleteRepo(repo *github.Repository) (wantsDelete bool, cancelled bool) {
	displayRepo(repo)
	prompt := promptui.Prompt{
		Label: fmt.Sprintf("Can Delete the repo: %s [use (Y/N)]", repo.GetName()),
	}

	result, err := prompt.Run()
	return promptDeleteRepoResult(result, err)
}

// shouldCallDelete is the one place that decides whether the GitHub API's
// delete call actually runs, kept pure and separate from SelectRepo's I/O so
// the dry-run/confirmation invariant can be tested without a terminal or a
// live GitHub client.
func shouldCallDelete(dryRun, isConfirmed bool) bool {
	return !dryRun && isConfirmed
}

func SelectRepo(login structures.Login, dry_run bool, repos []*github.Repository, classifier *bayesian.Classifier, forkFlag bool) {
	if classifier == nil {
		log.Println("Classifier not initialized, skipping sort.")
		return
	}

	sortedRepos := bayes.SortRepos(repos, classifier)
	nrepos := convertToGitHubRepos(sortedRepos)

	for _, repo := range nrepos {
		if !forkFlag && repo.GetFork() {
			continue
		}

		wantsDelete, cancelled := promptDeleteRepo(repo)
		if cancelled {
			fmt.Println("Skipping (prompt cancelled)...")
			continue
		}
		if !wantsDelete {
			// A deliberate "keep this repo" is as much a training signal as
			// a deletion, so record it instead of only ever learning Keep
			// from a delete confirmation answered "No".
			if err := files.SaveRepositoryFiles(repo, false); err != nil {
				log.Println("Failed to save repository file: ", err.Error())
			}
			fmt.Println("Skipping...")
			continue
		}

		isConfirmed := confirmDeletion(repo)
		if err := files.SaveRepositoryFiles(repo, isConfirmed); err != nil {
			log.Println("Failed to save repository file: ", err.Error())
		}

		if shouldCallDelete(dry_run, isConfirmed) {
			if err := gh.DeleteRepository(login, repo); err != nil {
				log.Println("Failed to delete repository: ", err.Error())
			}
		}
	}
}
