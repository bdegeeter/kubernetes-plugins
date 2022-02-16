//go:build integration
// +build integration

package integration

import (
	//"fmt"
	"os"
	"testing"

	"get.porter.sh/plugin/kubernetes/pkg/kubernetes/secrets"
	"get.porter.sh/plugin/kubernetes/tests"
	portercontext "get.porter.sh/porter/pkg/context"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"
)

var logger hclog.Logger = hclog.New(&hclog.LoggerOptions{
	Name:   secrets.PluginKey,
	Output: os.Stderr,
	Level:  hclog.Error})

func Test_Default_Namespace(t *testing.T) {
	k8sConfig := secrets.PluginConfig{}
	tc := portercontext.TestContext{}
	store := secrets.NewStore(tc.Context, k8sConfig)
	store.Connect()
	defer store.Close()
	t.Run("Test Default Namespace", func(t *testing.T) {
		_, err := store.Resolve("secret", "test")
		require.Error(t, err)
		if tests.RunningInKubernetes() {
			require.EqualError(t, err, "secrets \"test\" not found")
		} else {
			require.EqualError(t, err, "secrets \"test\" not found")
			//require.EqualError(t, err, "open /var/run/secrets/kubernetes.io/serviceaccount/namespace: no such file or directory")
		}
	})
}

func Test_Namespace_Does_Not_Exist(t *testing.T) {
	namespace := tests.GenerateNamespaceName()
	k8sConfig := secrets.PluginConfig{
		Namespace: namespace,
	}
	tc := portercontext.TestContext{}
	store := secrets.NewStore(tc.Context, k8sConfig)
	store.Connect()
	defer store.Close()
	t.Run("Test Namespace Does Not Exist", func(t *testing.T) {
		_, err := store.Resolve("secret", "test")
		require.Error(t, err)
		//TODO: sort out namespace error propagation
		//require.EqualError(t, err, fmt.Sprintf("namespaces \"%s\" not found", namespace))
	})
}

func TestResolve_Secret(t *testing.T) {
	nsName := tests.CreateNamespace(t)
	k8sConfig := secrets.PluginConfig{
		Namespace: nsName,
	}
	tc := portercontext.TestContext{}
	store := secrets.NewStore(tc.Context, k8sConfig)
	store.Connect()
	defer store.Close()
	defer tests.DeleteNamespace(t, nsName)
	tests.CreateSecret(t, nsName, "testkey", "testvalue")
	t.Run("resolve secret source: value", func(t *testing.T) {
		resolved, err := store.Resolve(secrets.SecretSourceType, "testkey")
		require.NoError(t, err)
		require.Equal(t, "testvalue", resolved)
	})

}

func Test_UppercaseKey(t *testing.T) {
	nsName := tests.CreateNamespace(t)
	defer tests.DeleteNamespace(t, nsName)
	k8sConfig := secrets.PluginConfig{
		Namespace: nsName,
	}
	tc := portercontext.TestContext{}
	store := secrets.NewStore(tc.Context, k8sConfig)
	store.Connect()
	defer store.Close()

	tests.CreateSecret(t, nsName, "testkey", "testvalue")
	t.Run("Test Uppercase Key", func(t *testing.T) {
		resolved, err := store.Resolve(secrets.SecretSourceType, "TESTkey")
		require.NoError(t, err)
		require.Equal(t, "testvalue", resolved)
	})
}
