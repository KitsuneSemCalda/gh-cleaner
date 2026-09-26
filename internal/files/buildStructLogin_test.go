package files

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeNetrc(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".netrc")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestMountLoginParsesGitHubEntry(t *testing.T) {
	path := writeNetrc(t, "machine github.com login someuser password sometoken\n")

	login, err := MountLogin(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if login.GetLogin() != "someuser" {
		t.Errorf("GetLogin() = %q, want %q", login.GetLogin(), "someuser")
	}
	if login.GetToken() != "sometoken" {
		t.Errorf("GetToken() = %q, want %q", login.GetToken(), "sometoken")
	}
}

// TestMountLoginErrorsWhenGitHubEntryMissing is the regression test for the
// audit finding: previously, a .netrc with no "machine github.com" entry
// (e.g. only unrelated machines, or an empty file) silently returned a
// zero-value Login and a nil error, which surfaced later as a confusing
// authentication failure against the GitHub API instead of a clear error at
// startup.
func TestMountLoginErrorsWhenGitHubEntryMissing(t *testing.T) {
	cases := map[string]string{
		"only unrelated machine": "machine example.com login someuser password sometoken\n",
		"empty file":             "",
	}

	for name, content := range cases {
		t.Run(name, func(t *testing.T) {
			path := writeNetrc(t, content)

			_, err := MountLogin(path)
			if err == nil {
				t.Fatal("expected an error when .netrc has no \"machine github.com\" entry, got nil")
			}
			if !strings.Contains(err.Error(), "github.com") {
				t.Errorf("error message %q does not mention the missing github.com entry", err.Error())
			}
		})
	}
}

func TestMountLoginErrorsWhenFileDoesNotExist(t *testing.T) {
	_, err := MountLogin(filepath.Join(t.TempDir(), "does-not-exist"))
	if err == nil {
		t.Fatal("expected an error for a missing .netrc file, got nil")
	}
}
