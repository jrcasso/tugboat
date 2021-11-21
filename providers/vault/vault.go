package vault

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/jrcasso/tugboat/tugboat"
)

var requiredEnvs = [...]string{
	"VAULT_ADDR",
	"VAULT_TOKEN",
}

type VaultProvider struct {
	Client  *api.Client
	Context *context.Context
}

var httpClient = &http.Client{
	Timeout: 15 * time.Second,
}

func CreateClient() *api.Client {
	tugboat.ValidateEnvironment(requiredEnvs[:])

	client, err := api.NewClient(&api.Config{Address: os.Getenv("VAULT_ADDR"), HttpClient: httpClient})
	if err != nil {
		panic(err)
	}
	client.SetToken(os.Getenv("VAULT_TOKEN"))

	return client
}

func Delete(service tugboat.Service) {
}

func Plan() []tugboat.ExecutionPlan {
	var plan []tugboat.ExecutionPlan
	// WHat secret do we needs?
	return plan
}

func (v VaultProvider) createSecret(key string, secret string) {
	inputData := map[string]interface{}{
		"data": map[string]interface{}{
			key: secret,
		},
	}
	output, err := v.Client.Logical().Write("secret/data/abd", inputData)
	fmt.Println(output)
	if err != nil {
		panic(err)
	}
}

func (v VaultProvider) listSecrets(dir string) {
	output, err := v.Client.Logical().Read(dir)
	fmt.Println(output)
	if err != nil {
		panic(err)
	}
}
