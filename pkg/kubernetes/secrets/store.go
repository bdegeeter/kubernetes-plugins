package secrets

import (
	"context"
	"fmt"
	"strings"

	"get.porter.sh/plugin/kubernetes/pkg/kubernetes/config"
	k8s "get.porter.sh/plugin/kubernetes/pkg/kubernetes/helper"
	portercontext "get.porter.sh/porter/pkg/context"
	portersecrets "get.porter.sh/porter/pkg/secrets/plugins"
	cnabsecrets "github.com/cnabio/cnab-go/secrets"
	"github.com/hashicorp/go-hclog"
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
	logger    hclog.Logger
	hostStore cnabsecrets.Store
	Secrets   map[string]map[string]string
	*portercontext.Context
	config    config.Config
	clientSet *kubernetes.Clientset
}

func NewStore(c *portercontext.Context, cfg PluginConfig) *Store {
	s := &Store{
		Secrets: make(map[string]map[string]string),
	}
	return s
}

func (s *Store) Connect() error {

	if s.clientSet != nil {
		return nil
	}

	clientSet, namespace, err := k8s.GetClientSet(s.config.Namespace, s.logger)

	if err != nil {
		s.logger.Debug(fmt.Sprintf("Failed to get Kubernetes Client Set: %v", err))
		return err
	}

	s.clientSet = clientSet
	s.config.Namespace = *namespace

	return nil
}

func (s *Store) Resolve(keyName string, keyValue string) (string, error) {
	if strings.ToLower(keyName) != SecretSourceType {
		return s.hostStore.Resolve(keyName, keyValue)
	}

	key := strings.ToLower(keyValue)

	s.logger.Debug(fmt.Sprintf("Looking for key:%s", keyValue))
	secret, err := s.clientSet.CoreV1().Secrets(s.config.Namespace).Get(context.Background(), key, metav1.GetOptions{})
	if err != nil {
		s.logger.Debug(fmt.Sprintf("Failed to Read secrets for key:%s %v", keyValue, err))
		return "", err
	}

	return string(secret.Data[SecretDataKey]), nil
}

func (s *Store) Close() error {

	/*
		if s.clientSet != nil {
			return nil
		}

		clientSet, namespace, err := k8s.GetClientSet(s.config.Namespace, s.logger)

		if err != nil {
			s.logger.Debug(fmt.Sprintf("Failed to get Kubernetes Client Set: %v", err))
			return err
		}

		s.clientSet = clientSet
		s.config.Namespace = *namespace
	*/

	return nil
}
