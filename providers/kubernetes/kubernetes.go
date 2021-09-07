package kubernetes

import (
	"context"
	"os"

	"github.com/jrcasso/tugboat/tugboat"
	log "github.com/sirupsen/logrus"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var requiredEnvs = [...]string{
	"KUBECONFIG",
}

type KubernetesProvider struct {
	Client  *kubernetes.Clientset
	Context *context.Context
}

func (k KubernetesProvider) Create(name string) {
	namespace := v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: name}}
	_, err := k.Client.CoreV1().Namespaces().Create(context.TODO(), &namespace, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Encountered error while creating namespace: %+v", err)
	}
	log.Infof("Successfully created namespace: %v", name)
}

func (k KubernetesProvider) Retrieve() []string {
	namespaceNames := []string{}
	namespaces, err := k.Client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Encountered error while listing namespaces: %+v", err)
	}
	for _, namespace := range namespaces.Items {
		log.Debugf("Found namespace: %+v", namespace.ObjectMeta.Name)
		namespaceNames = append(namespaceNames, namespace.ObjectMeta.Name)
	}

	return namespaceNames
}

func (k KubernetesProvider) Delete(name string) {
	err := k.Client.CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		log.Fatalf("Encountered error while listing namespaces: %+v", err)
	}
	log.Infof("Successfully deleted namespace: %v", name)
}

func (k KubernetesProvider) Execute(plan []tugboat.ExecutionPlan) {
	for _, command := range plan {
		command.Function(command.Arguments)
	}
}

// func (k KubernetesProvider) Plan(service tugboat.Service) []func(name string) {
func (k KubernetesProvider) Plan(services []tugboat.Service) []tugboat.ExecutionPlan {
	// executionPlan := []func(name string){}
	var namespaceExists bool
	executionPlan := []tugboat.ExecutionPlan{}
	namespaces := k.Retrieve()

	for _, service := range services {
		namespaceExists = false
		localNamespace := service.Name
		if service.Namespace != "" {
			// Allow user to override convention
			localNamespace = service.Namespace
		}

		for _, remoteNamespace := range namespaces {
			if localNamespace == remoteNamespace {
				namespaceExists = true
				break
			}
		}

		if !namespaceExists {
			executionPlan = append(
				executionPlan,
				tugboat.ExecutionPlan{
					Function:  k.Create,
					Arguments: localNamespace,
				})
		}
	}

	return executionPlan
}

func CreateClient() *kubernetes.Clientset {
	validateEnvironment()
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		log.Fatalf("Encountered error while parsing KUBECONFIG: %+v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Encountered error creating Kubernets client: %+v", err)
	}

	return clientset
}

func validateEnvironment() {
	for i := 0; i < len(requiredEnvs); i++ {
		var env = os.Getenv(requiredEnvs[i])
		if env == "" {
			log.Fatalf("Required environment variable not set: %+v", requiredEnvs[i])
		}
	}
}
