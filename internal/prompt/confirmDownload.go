package prompt

import (
	"fmt"
	"log"

	"github.com/google/go-github/v62/github"
	"github.com/manifoldco/promptui"
)

func confirmDownload(r *github.Repository) bool {
	prompt := promptui.Select{
		Label: "Do you really want delete the repo: " + r.GetName(),
		Items: []string{"Yes", "No"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalln(err)
	}

	switch result {
	case "Yes":
		fmt.Println("Return true")
		return true
	case "No":
		fmt.Println("Return false")
		return false
	default:
		fmt.Println("Return false")
		return false
	}
}
