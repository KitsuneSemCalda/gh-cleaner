package github

import (
	"context"
	"gh-cleaner/internal/structures"

	"github.com/google/go-github/v62/github"
)

// RepoDeleter is the thin slice of the GitHub API's Repositories service that
// DeleteRepository depends on. Defining it separately from *github.Client
// lets tests substitute a mock for gh-cleaner's one destructive, irreversible
// operation, instead of leaving it uncovered because it "needs a real
// GitHub client".
type RepoDeleter interface {
	Delete(ctx context.Context, owner, repo string) (*github.Response, error)
}

// DeleteRepository deletes the given repository from GitHub. This is the
// only irreversible action in gh-cleaner, so its actual API call is kept in
// the small, mockable deleteRepository below.
func DeleteRepository(l structures.Login, r *github.Repository) error {
	client := getClient(l)
	return deleteRepository(client.Repositories, l.GetLogin(), r.GetName())
}

// deleteRepository issues the real delete call through the RepoDeleter
// interface, isolated from client construction so it can be exercised with a
// mock in tests.
func deleteRepository(d RepoDeleter, owner, name string) error {
	_, err := d.Delete(context.Background(), owner, name)
	return err
}
