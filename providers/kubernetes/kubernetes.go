package kubernetes

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var requiredEnvs = [...]string{
	"KUBECONFIG",
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

func GetNamespaces(clientset kubernetes.Clientset) *v1.NamespaceList {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Encountered error while listing namespaces: %+v", err)
	}

	return namespaces
}

func validateEnvironment() {
	for i := 0; i < len(requiredEnvs); i++ {
		var env = os.Getenv(requiredEnvs[i])
		if env == "" {
			log.Fatalf("Required environment variable not set: %+v", requiredEnvs[i])
		}
	}
}
