package kubernetes

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/PapaDanielVi/secret-shift/internal/destination"
	"github.com/PapaDanielVi/secret-shift/internal/source"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Destination struct {
	clientset  *kubernetes.Clientset
	namespace  string
	secretName string
}

func New(kubeconfig, namespace, secretName string) (*Destination, error) {
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

	return &Destination{
		clientset:  clientset,
		namespace:  namespace,
		secretName: secretName,
	}, nil
}

func (d *Destination) Write(ctx context.Context, secrets []source.Secret) error {
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
		if err := d.writeSecret(ctx, secretData); err != nil {
			return fmt.Errorf("write secret: %w", err)
		}
	}

	if len(cmData) > 0 {
		if err := d.writeConfigMap(ctx, cmData); err != nil {
			return fmt.Errorf("write configmap: %w", err)
		}
	}

	return nil
}

func (d *Destination) writeSecret(ctx context.Context, data map[string]string) error {
	secretName := d.secretName
	if secretName == "" {
		secretName = "secret-shift-import"
	}

	existing, err := d.clientset.CoreV1().Secrets(d.namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("get existing secret: %w", err)
	}

	if errors.IsNotFound(err) {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: d.namespace,
			},
			StringData: data,
		}
		_, err = d.clientset.CoreV1().Secrets(d.namespace).Create(ctx, secret, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("create secret: %w", err)
		}
	} else {
		if existing.StringData == nil {
			existing.StringData = make(map[string]string)
		}
		for k, v := range data {
			existing.StringData[k] = v
		}
		_, err = d.clientset.CoreV1().Secrets(d.namespace).Update(ctx, existing, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("update secret: %w", err)
		}
	}

	return nil
}

func (d *Destination) writeConfigMap(ctx context.Context, data map[string]string) error {
	cmName := d.secretName
	if cmName == "" {
		cmName = "secret-shift-import"
	}

	existing, err := d.clientset.CoreV1().ConfigMaps(d.namespace).Get(ctx, cmName, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("get existing configmap: %w", err)
	}

	if errors.IsNotFound(err) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cmName,
				Namespace: d.namespace,
			},
			Data: data,
		}
		_, err = d.clientset.CoreV1().ConfigMaps(d.namespace).Create(ctx, cm, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("create configmap: %w", err)
		}
	} else {
		if existing.Data == nil {
			existing.Data = make(map[string]string)
		}
		for k, v := range data {
			existing.Data[k] = v
		}
		_, err = d.clientset.CoreV1().ConfigMaps(d.namespace).Update(ctx, existing, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("update configmap: %w", err)
		}
	}

	return nil
}

var _ destination.Destination = (*Destination)(nil)
