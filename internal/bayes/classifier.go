package bayes

import (
	"sort"

	"gh-cleaner/internal/structures"

	"github.com/google/go-github/v62/github"
	"github.com/jbrukh/bayesian"
)

const (
	Delete bayesian.Class = "Delete"
	Keep   bayesian.Class = "Keep"
)

func GenerateClassifier(deleteStuffs, keepStuffs []structures.RepoDatum) *bayesian.Classifier {
	classifier := bayesian.NewClassifier(Delete, Keep)

	trainClassifier(classifier, deleteStuffs, Delete)
	trainClassifier(classifier, keepStuffs, Keep)

	return classifier
}

func trainClassifier(classifier *bayesian.Classifier, samples []structures.RepoDatum, class bayesian.Class) {
	for _, sample := range samples {
		classifier.Learn(sample.Tokens(), class)
	}
}

func SortRepos(repos []*github.Repository, classifier *bayesian.Classifier) []*structures.Repository {
	classifiedRepos := make([]*structures.Repository, len(repos))

	for i, repo := range repos {
		r := structures.CreateRepository(repo)
		scores, _, _, _ := classifier.SafeProbScores(r.GetClassifierValues())
		if len(scores) > 0 {
			r.Probability = scores[0]
		}
		classifiedRepos[i] = &r
	}

	sort.Slice(classifiedRepos, func(i, j int) bool {
		return classifiedRepos[i].Probability > classifiedRepos[j].Probability
	})

	return classifiedRepos
}
