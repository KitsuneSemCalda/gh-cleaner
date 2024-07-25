package bayestheorem

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

// GenerateClassifier cria um classificador Bayesiano treinado com os dados fornecidos.
func GenerateClassifier(deleteStuffs, keepStuffs [][]string) *bayesian.Classifier {
	classifier := bayesian.NewClassifier(Delete, Keep)

	trainClassifier(classifier, deleteStuffs, Delete)
	trainClassifier(classifier, keepStuffs, Keep)

	return classifier
}

// trainClassifier treina o classificador com as amostras fornecidas para uma classe específica.
func trainClassifier(classifier *bayesian.Classifier, samples [][]string, class bayesian.Class) {
	for _, sample := range samples {
		classifier.Learn(sample, class)
	}
}

// classifyRepo classifica um repositório e retorna o escore de probabilidade.
func classifyRepo(repo *github.Repository, classifier *bayesian.Classifier) float64 {
	r := structures.CreateRepository(repo)
	scores, _, _ := classifier.LogScores(r.GetAllValues())

	var totalScore float64
	for _, score := range scores {
		totalScore += score
	}

	return totalScore
}

// SortRepos classifica e ordena os repositórios com base nas probabilidades calculadas.
func SortRepos(repos []*github.Repository, classifier *bayesian.Classifier) []*structures.Repository {
	classifiedRepos := make([]*structures.Repository, len(repos))

	for i, repo := range repos {
		r := structures.CreateRepository(repo)
		r.Probability = classifyRepo(repo, classifier)
		classifiedRepos[i] = &r
	}

	sort.Slice(classifiedRepos, func(i, j int) bool {
		return classifiedRepos[i].Probability > classifiedRepos[j].Probability
	})

	return classifiedRepos
}
