package github

import (
	"context"
	"gh-cleaner/internal/structures"

	"github.com/google/go-github/v62/github"
)

func GetRepositoriesByToken(l structures.Login) []*github.Repository {
	// var myRepos []*github.Repository
	client := github.NewClient(nil).WithAuthToken(l.GetToken())
	repos, _, _ := client.Repositories.ListByAuthenticatedUser(context.Background(), &github.RepositoryListByAuthenticatedUserOptions{
		Visibility:  "all",
		Direction:   "asc",
		Sort:        "created",
		Affiliation: "owner",
	})

	/**
			for _, repo := range repos {
				if !*repo.Fork {
					myRepos = append(myRepos, repo)
				}
			}

	  return myRepos
		**/

	return repos
}
