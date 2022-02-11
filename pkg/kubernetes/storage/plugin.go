package storage

import (
	"get.porter.sh/plugin/kubernetes/pkg/kubernetes/secrets"
	portercontext "get.porter.sh/porter/pkg/context"
	secretsplugin "get.porter.sh/porter/pkg/secrets/plugins"
	"get.porter.sh/porter/pkg/storage/plugins"
	"github.com/cnabio/cnab-go/utils/crud"
)

const PluginInterface = plugins.PluginInterface + ".kubernetes.storage"

var _ crud.Store = &Plugin{}

// Plugin is the plugin wrapper for accessing storage from Kubernetes Secrets.
// Secrets are used rather than config maps as there may be sensitive information in the data
type Plugin struct {
	crud.Store
}

func NewPlugin(cxt *portercontext.Context, pluginConfig interface{}) (secretsplugin.SecretsPlugin, error) {
	cfg := secrets.PluginConfig{}
	/*
		logger := hclog.New(&hclog.LoggerOptions{
			Name:       PluginInterface,
			Output:     os.Stderr,
			Level:      hclog.Debug,
			JSONFormat: true,
		})
	*/

	/*
		return &plugins.Plugin{
			Impl: &Plugin{
				Store: NewStore(cfg, logger),
			},
		}
	*/
	return secrets.NewStore(cxt, cfg), nil
}
