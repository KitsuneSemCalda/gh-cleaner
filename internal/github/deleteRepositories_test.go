package github

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-github/v62/github"
)

// mockRepoDeleter records how it was called so tests can assert exactly
// which repository would have been deleted, without ever touching the real
// GitHub API.
type mockRepoDeleter struct {
	calls   []struct{ owner, repo string }
	err     error
	callCnt int
}

func (m *mockRepoDeleter) Delete(ctx context.Context, owner, repo string) (*github.Response, error) {
	m.callCnt++
	m.calls = append(m.calls, struct{ owner, repo string }{owner, repo})
	return nil, m.err
}

func TestDeleteRepositoryCallsDeleteWithOwnerAndName(t *testing.T) {
	mock := &mockRepoDeleter{}

	err := deleteRepository(mock, "KitsuneSemCalda", "some-repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if mock.callCnt != 1 {
		t.Fatalf("expected exactly 1 call to Delete, got %d", mock.callCnt)
	}
	got := mock.calls[0]
	if got.owner != "KitsuneSemCalda" || got.repo != "some-repo" {
		t.Fatalf("Delete called with (%q, %q), want (%q, %q)", got.owner, got.repo, "KitsuneSemCalda", "some-repo")
	}
}

func TestDeleteRepositoryPropagatesError(t *testing.T) {
	wantErr := errors.New("403 Forbidden")
	mock := &mockRepoDeleter{err: wantErr}

	err := deleteRepository(mock, "owner", "repo")
	if !errors.Is(err, wantErr) {
		t.Fatalf("deleteRepository() error = %v, want %v", err, wantErr)
	}
}
