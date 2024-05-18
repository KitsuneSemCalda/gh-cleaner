package prompt

import (
	"fmt"
	gh "gh-cleaner/internal/github"
	"gh-cleaner/internal/structures"
	"log"
	"strings"

	"github.com/google/go-github/v62/github"
	"github.com/manifoldco/promptui"
)

func brokenRepo(repo *github.Repository) {
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

func SelectRepo(l structures.Login, r []*github.Repository, f bool) {
	var deletedRepos []*github.Repository
	for _, repo := range r {
		if f {
			fmt.Println()
			brokenRepo(repo)
			fmt.Println()

			prompt := promptui.Prompt{
				Label: fmt.Sprintf("Can Delete the repo: %s [use (Y/N)]", repo.GetName()),
			}

			result, err := prompt.Run()
			if err != nil {
				log.Fatalln(err)
			}

			lowered := strings.ToLower(result)

			switch lowered {
			case "y":
				deletedRepos = append(deletedRepos, repo)
			case "n":
				continue
			default:
				fmt.Println("Invalid input, skipping...")
			}
		} else {
			if !f && !*repo.Fork {
				fmt.Println()
				brokenRepo(repo)
				fmt.Println()

				prompt := promptui.Prompt{
					Label: fmt.Sprintf("Can Delete the repo: %s [use (Y/N)]", repo.GetName()),
				}

				result, err := prompt.Run()
				if err != nil {
					log.Fatalln(err)
				}

				lowered := strings.ToLower(result)

				switch lowered {
				case "y":
					deletedRepos = append(deletedRepos, repo)
				case "n":
					continue
				default:
					fmt.Println("Invalid input, skipping...")
				}
			}
		}
	}

	for _, fromDelete := range deletedRepos {
		isConfirmed := confirmDownload(fromDelete)
		gh.DeleteRepository(l, fromDelete, isConfirmed)
	}
}
