package github

import (
	"context"
	"gh-cleaner/internal/structures"

	"github.com/google/go-github/v62/github"
)

// RepoLister is the thin slice of the GitHub API's Repositories service that
// GetRepositoriesByToken depends on, defined separately from *github.Client
// so tests can substitute a mock instead of hitting the real API.
type RepoLister interface {
	ListByAuthenticatedUser(ctx context.Context, opts *github.RepositoryListByAuthenticatedUserOptions) ([]*github.Repository, *github.Response, error)
}

func GetRepositoriesByToken(l structures.Login) ([]*github.Repository, error) {
	client := getClient(l)
	return listRepositories(client.Repositories)
}

// listRepositories paginates through the authenticated user's repositories
// via the RepoLister interface, isolated from client construction so it can
// be exercised with a mock in tests.
func listRepositories(lister RepoLister) ([]*github.Repository, error) {
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
		repos, resp, err := lister.ListByAuthenticatedUser(ctx, opts)
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
