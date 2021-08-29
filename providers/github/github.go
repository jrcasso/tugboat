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

func GetRepository(ctx context.Context, client github.Client) {

}

func DeleteRepository(ctx context.Context, client github.Client, name string) {
	_, err := client.Repositories.Delete(ctx, os.Getenv("GITHUB_ORGANIZATION"), name)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Successfully deleted repo: %v\n", name)
}

func CreateClient(ctx context.Context) github.Client {
	validateEnvironment()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})
	tc := oauth2.NewClient(ctx, ts)
	return *github.NewClient(tc)
}

func CreateRepository(ctx context.Context, client github.Client, name string) {
	repo := &github.Repository{
		Name:        github.String(name),
		Private:     github.Bool(true),
		Description: github.String("Test repository"),
	}
	repo, _, err := client.Repositories.Create(ctx, os.Getenv("GITHUB_ORGANIZATION"), repo)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Successfully created new repo: %v\n", repo.GetName())
}

func validateEnvironment() {
	for i := 0; i < len(requiredEnvs); i++ {
		var env = os.Getenv(requiredEnvs[i])
		if env == "" {
			log.Fatalf("Required environment variable not set: %v", requiredEnvs[i])
		}
	}
}
