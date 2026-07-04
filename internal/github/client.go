package github

import (
	"gh-cleaner/internal/structures"

	"github.com/google/go-github/v62/github"
)

var sharedClient *github.Client

func getClient(l structures.Login) *github.Client {
	if sharedClient == nil {
		sharedClient = github.NewClient(nil).WithAuthToken(l.GetToken())
	}
	return sharedClient
}
