package secrets

import (
	"testing"

	portercontext "get.porter.sh/porter/pkg/context"
	"github.com/stretchr/testify/require"
)

func Test_NoNamespace(t *testing.T) {
	tc := portercontext.TestContext{}
	k8sConfig := PluginConfig{Namespace: "default"}
	store := NewStore(tc.Context, k8sConfig)
	store.Connect()
	defer store.Close()
	t.Run("Test No Namespace", func(t *testing.T) {
		_, err := store.Resolve("secret", "test")
		require.Error(t, err)
		//TODO: incluster or not incluster, that is the question
		//require.EqualError(t, err, "open /var/run/secrets/kubernetes.io/serviceaccount/namespace: no such file or directory")
		require.EqualError(t, err, "secrets \"test\" not found")
	})
}
