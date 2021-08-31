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

func (k KubernetesProvider) Retrieve() interface{} {
	namespaces, err := k.Client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Encountered error while listing namespaces: %+v", err)
	}
	for _, namespace := range namespaces.Items {
		log.Debugf("Found namespace: %+v", namespace.ObjectMeta.Name)
	}

	return namespaces
}

func (k KubernetesProvider) Delete(name string) {
	err := k.Client.CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		log.Fatalf("Encountered error while listing namespaces: %+v", err)
	}
	log.Infof("Successfully deleted namespace: %v", name)
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
