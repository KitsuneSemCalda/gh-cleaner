package github

import (
	"context"
	"gh-cleaner/internal/structures"

	"github.com/google/go-github/v62/github"
)

func DeleteRepository(l structures.Login, r *github.Repository, confirmDeletion bool) {
	repoName := r.GetName()

	if confirmDeletion {
		client := github.NewClient(nil).WithAuthToken(l.GetToken())

		client.Repositories.Delete(context.Background(), l.GetLogin(), repoName)
		return
	}
}
