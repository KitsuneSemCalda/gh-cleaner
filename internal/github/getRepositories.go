package github

import (
	"context"
	"gh-cleaner/internal/structures"

	"github.com/google/go-github/v62/github"
)

func GetRepositoriesByToken(l structures.Login) ([]*github.Repository, error) {
	client := getClient(l)
	ctx := context.Background()

	var allRepos []*github.Repository
	opts := &github.RepositoryListByAuthenticatedUserOptions{
		Visibility:  "all",
		Direction:   "asc",
		Sort:        "created",
		Affiliation: "owner",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		repos, resp, err := client.Repositories.ListByAuthenticatedUser(ctx, opts)
		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repos...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allRepos, nil
}
