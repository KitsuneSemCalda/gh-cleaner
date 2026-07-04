package structures

import (
	"strconv"

	"github.com/google/go-github/v62/github"
)

type Repository struct {
	RepositoryName        string
	RepositoryStar        string
	RepositoryIssues      string
	RepositoryForks       string
	RepositorySubscribers string
	Probability           float64

	OriginalRepo *github.Repository
}

func CreateRepository(g *github.Repository) Repository {
	stars := strconv.Itoa(g.GetStargazersCount())
	issues := strconv.Itoa(g.GetOpenIssues())
	forks := strconv.Itoa(g.GetForksCount())
	subscribers := strconv.Itoa(g.GetSubscribersCount())

	return Repository{
		RepositoryName:        g.GetName(),
		RepositoryStar:        stars,
		RepositoryIssues:      issues,
		RepositoryForks:       forks,
		RepositorySubscribers: subscribers,

		Probability: 0.0,

		OriginalRepo: g,
	}
}

func (r *Repository) DataFields() []string {
	return []string{
		r.RepositoryName,
		r.RepositoryStar,
		r.RepositoryIssues,
		r.RepositoryForks,
		r.RepositorySubscribers,
	}
}

func (r *Repository) GetClassifierValues() []string {
	return []string{
		r.RepositoryStar,
		r.RepositoryIssues,
		r.RepositoryForks,
		r.RepositorySubscribers,
	}
}

type RepoDatum struct {
	Name        string
	Stars       string
	Issues      string
	Forks       string
	Subscribers string
}

func (rd RepoDatum) Tokens() []string {
	return []string{rd.Stars, rd.Issues, rd.Forks, rd.Subscribers}
}
