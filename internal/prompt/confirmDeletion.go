package prompt

import (
	"github.com/google/go-github/v62/github"
	"github.com/manifoldco/promptui"
)

// confirmDeletionResult interprets a promptui.Select outcome for the delete
// confirmation. Kept separate from confirmDeletion so it can be tested
// without a terminal: a cancelled or otherwise failed prompt must never be
// read as consent to delete.
func confirmDeletionResult(result string, err error) bool {
	if err != nil {
		return false
	}
	return result == "Yes"
}

func confirmDeletion(r *github.Repository) bool {
	prompt := promptui.Select{
		Label: "Do you really want delete the repo: " + r.GetName(),
		// "No" first so it is the item pre-selected when the prompt opens;
		// deleting a repo should need a deliberate move to "Yes", not Enter.
		Items: []string{"No", "Yes"},
	}

	_, result, err := prompt.Run()
	return confirmDeletionResult(result, err)
}
