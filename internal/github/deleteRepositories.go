package github

import (
	"context"
	"gh-cleaner/internal/structures"

	"github.com/google/go-github/v62/github"
)

func DeleteRepository(l structures.Login, r *github.Repository) error {
	client := getClient(l)
	_, err := client.Repositories.Delete(context.Background(), l.GetLogin(), r.GetName())
	return err
}
