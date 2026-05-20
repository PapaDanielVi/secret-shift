package kubernetes

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/PapaDanielVi/secret-shift/internal/source"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Source struct {
	clientset    *kubernetes.Clientset
	namespace    string
	secretName   string
	labelSelector string
}

func New(kubeconfig, namespace, secretName, labelSelector string) (*Source, error) {
	if kubeconfig == "" {
		kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("build kubeconfig from %s: %w", kubeconfig, err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("create kubernetes client: %w", err)
	}

	return &Source{
		clientset:     clientset,
		namespace:     namespace,
		secretName:    secretName,
		labelSelector: labelSelector,
	}, nil
}

func (s *Source) Read(ctx context.Context) ([]source.Secret, error) {
	var result []source.Secret

	if s.secretName != "" {
		secret, err := s.clientset.CoreV1().Secrets(s.namespace).Get(ctx, s.secretName, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("get secret %s/%s: %w", s.namespace, s.secretName, err)
		}
		for k, v := range secret.Data {
			result = append(result, source.Secret{
				Name:  k,
				Value: string(v),
				Type:  "secret",
			})
		}
		return result, nil
	}

	if s.labelSelector != "" {
		secrets, err := s.clientset.CoreV1().Secrets(s.namespace).List(ctx, metav1.ListOptions{
			LabelSelector: s.labelSelector,
		})
		if err != nil {
			return nil, fmt.Errorf("list secrets with label %s: %w", s.labelSelector, err)
		}
		for _, secret := range secrets.Items {
			for k, v := range secret.Data {
				name := secret.Name + "/" + k
				result = append(result, source.Secret{
					Name:  name,
					Value: string(v),
					Type:  "secret",
				})
			}
		}
		return result, nil
	}

	// Read all secrets in namespace
	secrets, err := s.clientset.CoreV1().Secrets(s.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list secrets in namespace %s: %w", s.namespace, err)
	}
	for _, secret := range secrets.Items {
		for k, v := range secret.Data {
			name := secret.Name + "/" + k
			result = append(result, source.Secret{
				Name:  name,
				Value: string(v),
				Type:  "secret",
			})
		}
	}

	// Also read configmaps
	configmaps, err := s.clientset.CoreV1().ConfigMaps(s.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list configmaps in namespace %s: %w", s.namespace, err)
	}
	for _, cm := range configmaps.Items {
		for k, v := range cm.Data {
			name := cm.Name + "/" + k
			result = append(result, source.Secret{
				Name:  name,
				Value: v,
				Type:  "env",
			})
		}
	}

	return result, nil
}
