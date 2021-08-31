package github

import (
	"context"
	"os"

	"github.com/google/go-github/github"
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

func (g GithubProvider) Retrieve() interface{} {
	repos, _, err := g.Client.Repositories.ListByOrg(*g.Context, os.Getenv("GITHUB_ORGANIZATION"), nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, repo := range repos {
		log.Debugf("Found repo: %+v", *repo.Name)
	}

	return repos
}

func (g GithubProvider) Delete(name string) {
	_, err := g.Client.Repositories.Delete(*g.Context, os.Getenv("GITHUB_ORGANIZATION"), name)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Successfully deleted repo: %v", name)
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
