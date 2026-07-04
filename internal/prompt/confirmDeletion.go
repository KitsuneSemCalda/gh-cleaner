package prompt

import (
	"github.com/google/go-github/v62/github"
	"github.com/manifoldco/promptui"
)

func confirmDeletion(r *github.Repository) bool {
	prompt := promptui.Select{
		Label: "Do you really want delete the repo: " + r.GetName(),
		Items: []string{"Yes", "No"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return false
	}

	return result == "Yes"
}
