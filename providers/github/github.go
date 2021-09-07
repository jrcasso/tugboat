package github

import (
	"context"
	"os"

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

func (g GithubProvider) Create(name string) {
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

func (g GithubProvider) Retrieve() []string {
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

func (g GithubProvider) Delete(name string) {
	_, err := g.Client.Repositories.Delete(*g.Context, os.Getenv("GITHUB_ORGANIZATION"), name)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Successfully deleted repo: %v", name)
}

func (g GithubProvider) Execute(plan []tugboat.ExecutionPlan) {
	for _, command := range plan {
		command.Function(command.Arguments)
	}
}

func (g GithubProvider) Plan(services []tugboat.Service) []tugboat.ExecutionPlan {
	var repoExists bool
	executionPlan := []tugboat.ExecutionPlan{}
	repos := g.Retrieve()

	for _, service := range services {
		repoExists = false
		localRepo := service.Name
		if service.Repo != "" {
			// Allow user to override convention
			localRepo = service.Repo
		}

		for _, remoteRepo := range repos {
			if localRepo == remoteRepo {
				repoExists = true
				break
			}
		}

		if !repoExists {
			executionPlan = append(
				executionPlan,
				tugboat.ExecutionPlan{
					Function:  g.Create,
					Arguments: localRepo,
				})
		}
	}

	return executionPlan
}

func CreateClient(ctx context.Context) github.Client {
	validateEnvironment()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})
	tc := oauth2.NewClient(ctx, ts)

	return *github.NewClient(tc)
}

func validateEnvironment() {
	for _, env := range requiredEnvs {
		var env = os.Getenv(env)
		if env == "" {
			log.Fatalf("Required environment variable not set: %v", env)
		}
	}
}
