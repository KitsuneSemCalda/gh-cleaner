package structures

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/google/go-github/v62/github"
)

type Repository struct {
	RepositoryName   string
	RepositoryStar   string
	RepositoryIssues string
	RepositoryForks  string
}

func CreateRepository(g *github.Repository) Repository {
	stars := strconv.Itoa(g.GetStargazersCount())
	issues := strconv.Itoa(g.GetOpenIssues())
	forks := strconv.Itoa(g.GetOpenIssues())

	return Repository{
		RepositoryName:   g.GetName(),
		RepositoryStar:   stars,
		RepositoryIssues: issues,
		RepositoryForks:  forks,
	}
}

func (r *Repository) GetAllValues() []string {
	var fieldNames []string

	val := reflect.ValueOf(r)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < val.NumField(); i++ {
		fieldValue := val.Field(i).Interface()
		fieldNames = append(fieldNames, fmt.Sprintf("%v", fieldValue))
	}

	return fieldNames
}
