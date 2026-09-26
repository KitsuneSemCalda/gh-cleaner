package prompt

import (
	"errors"
	"testing"

	"github.com/google/go-github/v62/github"
)

func TestPromptDeleteRepoResult(t *testing.T) {
	cases := []struct {
		name          string
		result        string
		err           error
		wantsDelete   bool
		wantCancelled bool
	}{
		{"lowercase y wants delete", "y", nil, true, false},
		{"uppercase Y wants delete", "Y", nil, true, false},
		{"n keeps the repo", "n", nil, false, false},
		{"anything else keeps the repo", "whatever", nil, false, false},
		{"empty input keeps the repo", "", nil, false, false},
		{"cancelled prompt is not a keep decision", "", errors.New("^C"), false, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotDelete, gotCancelled := promptDeleteRepoResult(c.result, c.err)
			if gotDelete != c.wantsDelete || gotCancelled != c.wantCancelled {
				t.Errorf("promptDeleteRepoResult(%q, %v) = (%v, %v), want (%v, %v)",
					c.result, c.err, gotDelete, gotCancelled, c.wantsDelete, c.wantCancelled)
			}
		})
	}
}

func TestShouldCallDelete(t *testing.T) {
	cases := []struct {
		dryRun      bool
		isConfirmed bool
		want        bool
	}{
		{dryRun: true, isConfirmed: true, want: false},
		{dryRun: true, isConfirmed: false, want: false},
		{dryRun: false, isConfirmed: false, want: false},
		{dryRun: false, isConfirmed: true, want: true},
	}
	for _, c := range cases {
		if got := shouldCallDelete(c.dryRun, c.isConfirmed); got != c.want {
			t.Errorf("shouldCallDelete(%v, %v) = %v, want %v", c.dryRun, c.isConfirmed, got, c.want)
		}
	}
}

// TestProcessRepo_DryRunBlocksRealDelete is the end-to-end regression test
// for gh-cleaner's single destructive operation: with --dry-run set, even a
// fully confirmed delete must never reach the (here, faked) GitHub API call.
// This extends the pure shouldCallDelete coverage above by wiring the actual
// SelectRepo code path (via processRepo) against a spy deleteFn.
func TestProcessRepo_DryRunBlocksRealDelete(t *testing.T) {
	repo := &github.Repository{Name: github.String("test-repo")}
	deleteCalls := 0
	deleteFn := func(r *github.Repository) error {
		deleteCalls++
		return nil
	}
	saveFn := func(r *github.Repository, deleted bool) error { return nil }

	processRepo(repo, true /* dryRun */, true /* wantsDelete */, false /* cancelled */, true /* isConfirmed */, deleteFn, saveFn)

	if deleteCalls != 0 {
		t.Fatalf("dry-run must block the real delete call, but deleteFn was invoked %d time(s)", deleteCalls)
	}
}

func TestProcessRepo_ConfirmedNonDryRunCallsDelete(t *testing.T) {
	repo := &github.Repository{Name: github.String("test-repo")}
	deleteCalls := 0
	var deletedRepo *github.Repository
	deleteFn := func(r *github.Repository) error {
		deleteCalls++
		deletedRepo = r
		return nil
	}
	saveFn := func(r *github.Repository, deleted bool) error { return nil }

	processRepo(repo, false /* dryRun */, true /* wantsDelete */, false /* cancelled */, true /* isConfirmed */, deleteFn, saveFn)

	if deleteCalls != 1 {
		t.Fatalf("expected exactly 1 delete call for a confirmed, non-dry-run deletion, got %d", deleteCalls)
	}
	if deletedRepo != repo {
		t.Fatalf("deleteFn was called with the wrong repository")
	}
}

func TestProcessRepo_NotConfirmedNeverCallsDeleteEvenWithoutDryRun(t *testing.T) {
	repo := &github.Repository{Name: github.String("test-repo")}
	deleteCalls := 0
	deleteFn := func(r *github.Repository) error {
		deleteCalls++
		return nil
	}
	saveFn := func(r *github.Repository, deleted bool) error { return nil }

	processRepo(repo, false /* dryRun */, true /* wantsDelete */, false /* cancelled */, false /* isConfirmed */, deleteFn, saveFn)

	if deleteCalls != 0 {
		t.Fatalf("a declined confirmation must never trigger delete, but deleteFn was invoked %d time(s)", deleteCalls)
	}
}

func TestProcessRepo_CancelledPromptSkipsWithoutSavingOrDeleting(t *testing.T) {
	repo := &github.Repository{Name: github.String("test-repo")}
	deleteCalls, saveCalls := 0, 0
	deleteFn := func(r *github.Repository) error {
		deleteCalls++
		return nil
	}
	saveFn := func(r *github.Repository, deleted bool) error {
		saveCalls++
		return nil
	}

	processRepo(repo, false /* dryRun */, false /* wantsDelete */, true /* cancelled */, false /* isConfirmed */, deleteFn, saveFn)

	if deleteCalls != 0 || saveCalls != 0 {
		t.Fatalf("a cancelled prompt must skip without saving or deleting, got deleteCalls=%d saveCalls=%d", deleteCalls, saveCalls)
	}
}

func TestProcessRepo_DeclinedWantsDeleteSavesKeepDecision(t *testing.T) {
	repo := &github.Repository{Name: github.String("test-repo")}
	var savedDeleted *bool
	deleteFn := func(r *github.Repository) error { return nil }
	saveFn := func(r *github.Repository, deleted bool) error {
		savedDeleted = &deleted
		return nil
	}

	processRepo(repo, false /* dryRun */, false /* wantsDelete */, false /* cancelled */, false /* isConfirmed */, deleteFn, saveFn)

	if savedDeleted == nil {
		t.Fatal("expected saveFn to be called recording the keep decision")
	}
	if *savedDeleted != false {
		t.Fatalf("expected saveFn to be called with deleted=false, got %v", *savedDeleted)
	}
}
