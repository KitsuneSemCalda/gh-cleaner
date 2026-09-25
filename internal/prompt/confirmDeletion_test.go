package prompt

import (
	"errors"
	"testing"
)

func TestConfirmDeletionResult(t *testing.T) {
	cases := []struct {
		name   string
		result string
		err    error
		want   bool
	}{
		{"explicit yes", "Yes", nil, true},
		{"explicit no", "No", nil, false},
		{"cancelled prompt never confirms", "", errors.New("^C"), false},
		{"prompt error with stray result never confirms", "Yes", errors.New("io error"), false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := confirmDeletionResult(c.result, c.err); got != c.want {
				t.Errorf("confirmDeletionResult(%q, %v) = %v, want %v", c.result, c.err, got, c.want)
			}
		})
	}
}
