package prompt

import (
	"errors"
	"testing"
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
