package kubernetes

import (
	"context"
	"fmt"
	"maps"
	"os"
	"path/filepath"

	"github.com/PapaDanielVi/secret-shift/internal/provider"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Provider struct {
	clientset     *kubernetes.Clientset
	namespace     string
	secretName    string
	labelSelector string
}

func New(kubeconfig, namespace, secretName, labelSelector string) (*Provider, error) {
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

	return &Provider{
		clientset:     clientset,
		namespace:     namespace,
		secretName:    secretName,
		labelSelector: labelSelector,
	}, nil
}

func (p *Provider) Read(ctx context.Context) ([]provider.Secret, error) {
	if p.secretName != "" {
		return p.readSingleSecret(ctx)
	}

	if p.labelSelector != "" {
		return p.readSecretsByLabel(ctx)
	}

	return p.readAllSecretsAndConfigMaps(ctx)
}

func (p *Provider) readSingleSecret(ctx context.Context) ([]provider.Secret, error) {
	secret, err := p.clientset.CoreV1().Secrets(p.namespace).Get(ctx, p.secretName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("get secret %s/%s: %w", p.namespace, p.secretName, err)
	}
	return p.extractSecretsFromSecret(secret, ""), nil
}

func (p *Provider) readSecretsByLabel(ctx context.Context) ([]provider.Secret, error) {
	secrets, err := p.clientset.CoreV1().Secrets(p.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: p.labelSelector,
	})
	if err != nil {
		return nil, fmt.Errorf("list secrets with label %s: %w", p.labelSelector, err)
	}

	var result []provider.Secret
	for _, secret := range secrets.Items {
		result = append(result, p.extractSecretsFromSecret(&secret, secret.Name+"/")...)
	}
	return result, nil
}

func (p *Provider) readAllSecretsAndConfigMaps(ctx context.Context) ([]provider.Secret, error) {
	var result []provider.Secret

	secrets, err := p.clientset.CoreV1().Secrets(p.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list secrets in namespace %s: %w", p.namespace, err)
	}
	for _, secret := range secrets.Items {
		result = append(result, p.extractSecretsFromSecret(&secret, secret.Name+"/")...)
	}

	configmaps, err := p.clientset.CoreV1().ConfigMaps(p.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list configmaps in namespace %s: %w", p.namespace, err)
	}
	for _, cm := range configmaps.Items {
		for k, v := range cm.Data {
			result = append(result, provider.Secret{
				Name:  cm.Name + "/" + k,
				Value: v,
				Type:  "env",
			})
		}
	}

	return result, nil
}

func (p *Provider) extractSecretsFromSecret(secret *corev1.Secret, prefix string) []provider.Secret {
	var result []provider.Secret
	for k, v := range secret.Data {
		result = append(result, provider.Secret{
			Name:  prefix + k,
			Value: string(v),
			Type:  "secret",
		})
	}
	return result
}

func (p *Provider) Write(ctx context.Context, secrets []provider.Secret) error {
	secretData := make(map[string]string)
	cmData := make(map[string]string)

	for _, s := range secrets {
		if s.Type == "secret" {
			secretData[s.Name] = s.Value
		} else {
			cmData[s.Name] = s.Value
		}
	}

	if len(secretData) > 0 {
		if err := p.writeSecret(ctx, secretData); err != nil {
			return fmt.Errorf("write secret: %w", err)
		}
	}

	if len(cmData) > 0 {
		if err := p.writeConfigMap(ctx, cmData); err != nil {
			return fmt.Errorf("write configmap: %w", err)
		}
	}

	return nil
}

func (p *Provider) writeSecret(ctx context.Context, data map[string]string) error { //nolint:dupl
	name := p.secretName
	if name == "" {
		name = "secret-shift-import"
	}

	existing, err := p.clientset.CoreV1().Secrets(p.namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("get existing secret: %w", err)
	}

	if errors.IsNotFound(err) {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: p.namespace,
			},
			StringData: data,
		}
		_, err = p.clientset.CoreV1().Secrets(p.namespace).Create(ctx, secret, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("create secret: %w", err)
		}
	} else {
		if existing.StringData == nil {
			existing.StringData = make(map[string]string)
		}
		maps.Insert(existing.StringData, maps.All(data))
		_, err = p.clientset.CoreV1().Secrets(p.namespace).Update(ctx, existing, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("update secret: %w", err)
		}
	}

	return nil
}

func (p *Provider) writeConfigMap(ctx context.Context, data map[string]string) error { //nolint:dupl
	name := p.secretName
	if name == "" {
		name = "secret-shift-import"
	}

	existing, err := p.clientset.CoreV1().ConfigMaps(p.namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("get existing configmap: %w", err)
	}

	if errors.IsNotFound(err) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: p.namespace,
			},
			Data: data,
		}
		_, err = p.clientset.CoreV1().ConfigMaps(p.namespace).Create(ctx, cm, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("create configmap: %w", err)
		}
	} else {
		if existing.Data == nil {
			existing.Data = make(map[string]string)
		}
		maps.Insert(existing.Data, maps.All(data))
		_, err = p.clientset.CoreV1().ConfigMaps(p.namespace).Update(ctx, existing, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("update configmap: %w", err)
		}
	}

	return nil
}

var _ provider.Source = (*Provider)(nil)
var _ provider.Destination = (*Provider)(nil)
