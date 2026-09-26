package github

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-github/v62/github"
)

// mockRepoLister simulates a paginated ListByAuthenticatedUser without
// hitting the real GitHub API.
type mockRepoLister struct {
	pages [][]*github.Repository
	err   error
	calls int
}

func (m *mockRepoLister) ListByAuthenticatedUser(ctx context.Context, opts *github.RepositoryListByAuthenticatedUserOptions) ([]*github.Repository, *github.Response, error) {
	m.calls++
	if m.err != nil {
		return nil, nil, m.err
	}

	page := opts.ListOptions.Page
	if page == 0 {
		page = 1
	}
	repos := m.pages[page-1]

	resp := &github.Response{}
	if page < len(m.pages) {
		resp.NextPage = page + 1
	}
	return repos, resp, nil
}

func TestListRepositoriesAggregatesAllPages(t *testing.T) {
	repoA := &github.Repository{Name: github.String("repo-a")}
	repoB := &github.Repository{Name: github.String("repo-b")}
	repoC := &github.Repository{Name: github.String("repo-c")}

	mock := &mockRepoLister{pages: [][]*github.Repository{
		{repoA, repoB},
		{repoC},
	}}

	repos, err := listRepositories(mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.calls != 2 {
		t.Fatalf("expected 2 pages fetched, got %d calls", mock.calls)
	}
	if len(repos) != 3 {
		t.Fatalf("expected 3 aggregated repos, got %d", len(repos))
	}
	names := []string{repos[0].GetName(), repos[1].GetName(), repos[2].GetName()}
	want := []string{"repo-a", "repo-b", "repo-c"}
	for i := range want {
		if names[i] != want[i] {
			t.Errorf("repos[%d] = %q, want %q", i, names[i], want[i])
		}
	}
}

func TestListRepositoriesPropagatesError(t *testing.T) {
	wantErr := errors.New("rate limited")
	mock := &mockRepoLister{err: wantErr}

	_, err := listRepositories(mock)
	if !errors.Is(err, wantErr) {
		t.Fatalf("listRepositories() error = %v, want %v", err, wantErr)
	}
}
