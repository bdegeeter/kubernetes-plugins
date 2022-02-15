package secrets

import (
	"context"
	"fmt"
	"strings"

	k8s "get.porter.sh/plugin/kubernetes/pkg/kubernetes/helper"
	portercontext "get.porter.sh/porter/pkg/context"
	portersecrets "get.porter.sh/porter/pkg/secrets/plugins"
	cnabsecrets "github.com/cnabio/cnab-go/secrets"
	cnabhost "github.com/cnabio/cnab-go/secrets/host"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var _ portersecrets.SecretsPlugin = &Store{}

const (
	SecretSourceType = "secret"
	SecretDataKey    = "credential"
)

// Store implements the backing store for secrets as kubernetes secrets.
type Store struct {
	*portercontext.Context
	hostStore  cnabsecrets.Store
	Secrets    map[string]map[string]string
	namespace  string
	kubeconfig string
	clientSet  *kubernetes.Clientset
}

func NewStore(c *portercontext.Context, cfg PluginConfig) *Store {
	namespace := cfg.Namespace
	s := &Store{
		Secrets:   make(map[string]map[string]string),
		hostStore: &cnabhost.SecretStore{},
		namespace: namespace,
	}
	return s
}

func (s *Store) Connect() error {

	if s.clientSet != nil {
		return nil
	}
	clientSet, namespace, err := k8s.GetClientSet(s.namespace)

	if err != nil {
		return err
	}

	s.clientSet = clientSet
	s.namespace = *namespace

	return nil
}

func (s *Store) Resolve(keyName string, keyValue string) (string, error) {
	if strings.ToLower(keyName) != SecretSourceType {
		return s.hostStore.Resolve(keyName, keyValue)
	}

	key := strings.ToLower(keyValue)

	secret, err := s.clientSet.CoreV1().Secrets(s.namespace).Get(context.Background(), key, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("\nFailed to Read secrets for key: %s %v\n", keyValue, err)
		return "", err
	}

	return string(secret.Data[SecretDataKey]), nil
}

func (s *Store) Close() error {

	if s.clientSet != nil {
		return nil
	}

	clientSet, namespace, err := k8s.GetClientSet(s.namespace)

	if err != nil {
		return err
	}

	s.clientSet = clientSet
	s.namespace = *namespace

	return nil
}
