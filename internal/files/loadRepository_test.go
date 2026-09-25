package files

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFixtureRepo(t *testing.T, dir, name string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	content := name + "\n1\n2\n3\n4\n"
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestGetInfoAboutRepoHandlesMissingKeepOnlyDeleteOnlyAndMixedHistories(t *testing.T) {
	t.Run("neither directory exists yet (first run)", func(t *testing.T) {
		t.Setenv("HOME", t.TempDir())
		saved, deleted, err := GetInfoAboutRepo()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(saved) != 0 || len(deleted) != 0 {
			t.Fatalf("expected empty classes, got saved=%v deleted=%v", saved, deleted)
		}
	})

	t.Run("only Keep decisions recorded so far", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)
		writeFixtureRepo(t, filepath.Join(home, ".local", "share", "gh-cleaner", "repository", "saved"), "keep-only")

		saved, deleted, err := GetInfoAboutRepo()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(saved) != 1 || saved[0].Name != "keep-only" {
			t.Fatalf("expected one saved repo, got %v", saved)
		}
		if len(deleted) != 0 {
			t.Fatalf("expected no deleted repos, got %v", deleted)
		}
	})

	t.Run("only Delete decisions recorded so far", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)
		writeFixtureRepo(t, filepath.Join(home, ".local", "share", "gh-cleaner", "repository", "deleted"), "delete-only")

		saved, deleted, err := GetInfoAboutRepo()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(deleted) != 1 || deleted[0].Name != "delete-only" {
			t.Fatalf("expected one deleted repo, got %v", deleted)
		}
		if len(saved) != 0 {
			t.Fatalf("expected no saved repos, got %v", saved)
		}
	})

	t.Run("mixed history in both directories", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)
		base := filepath.Join(home, ".local", "share", "gh-cleaner", "repository")
		writeFixtureRepo(t, filepath.Join(base, "saved"), "kept-repo")
		writeFixtureRepo(t, filepath.Join(base, "deleted"), "deleted-repo")

		saved, deleted, err := GetInfoAboutRepo()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(saved) != 1 || saved[0].Name != "kept-repo" {
			t.Fatalf("expected one saved repo, got %v", saved)
		}
		if len(deleted) != 1 || deleted[0].Name != "deleted-repo" {
			t.Fatalf("expected one deleted repo, got %v", deleted)
		}
	})
}
