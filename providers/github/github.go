package github

import (
	"context"
	"os"
	"sync"

	"github.com/google/go-github/github"
	"github.com/jrcasso/tugboat/tugboat"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

var requiredEnvs = [...]string{
	"GITHUB_TOKEN",
	"GITHUB_ORGANIZATION",
}

type GithubProvider struct {
	Client  *github.Client
	Context *context.Context
}

func (g GithubProvider) create(name string) {
	repo := &github.Repository{
		Name:        github.String(name),
		Private:     github.Bool(true),
		Description: github.String("Test repository"),
	}
	repo, _, err := g.Client.Repositories.Create(*g.Context, os.Getenv("GITHUB_ORGANIZATION"), repo)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Successfully created repo: %v", repo.GetName())
}

func (g GithubProvider) retrieve() []string {
	repoNames := []string{}
	repos, _, err := g.Client.Repositories.ListByOrg(*g.Context, os.Getenv("GITHUB_ORGANIZATION"), nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, repo := range repos {
		log.Debugf("Found repo: %+v", *repo.Name)
		repoNames = append(repoNames, *repo.Name)
	}

	return repoNames
}

func (g GithubProvider) Delete(service tugboat.Service) {
	_, err := g.Client.Repositories.Delete(*g.Context, os.Getenv("GITHUB_ORGANIZATION"), service.Name)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Successfully deleted repo: %v", service.Name)
}

func (g GithubProvider) Plan(services []tugboat.Service, wg *sync.WaitGroup) []tugboat.ExecutionPlan {
	defer wg.Done()
	var repoExists bool
	executionPlan := []tugboat.ExecutionPlan{}
	repos := g.retrieve()

	for _, service := range services {
		repoExists = false
		localRepo := service.Name
		if service.Repo != "" {
			// Allow user to override convention
			localRepo = service.Repo
		}

		if tugboat.SliceContains(localRepo, repos) {
			repoExists = true
		}

		if !repoExists {
			executionPlan = append(
				executionPlan,
				tugboat.ExecutionPlan{
					Function:  g.create,
					Arguments: localRepo,
				})
		}
	}

	return executionPlan
}

func CreateClient(ctx context.Context) github.Client {
	tugboat.ValidateEnvironment(requiredEnvs[:])
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})
	tc := oauth2.NewClient(ctx, ts)

	return *github.NewClient(tc)
}
