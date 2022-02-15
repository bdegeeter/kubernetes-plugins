package secrets

import (
	portercontext "get.porter.sh/porter/pkg/context"
	"get.porter.sh/porter/pkg/secrets/plugins"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

const PluginKey = plugins.PluginInterface + ".kubernetes.secret"

type PluginConfig struct {
	KubeConfig string `mapstructure:"kubeconfig"`
	Namespace  string `mapstructure:"namespace"`
}

func NewPlugin(cxt *portercontext.Context, pluginConfig interface{}) (plugins.SecretsPlugin, error) {
	cfg := PluginConfig{}

	if err := mapstructure.Decode(pluginConfig, &cfg); err != nil {
		return nil, errors.Wrapf(err, "error decoding %s plugin config from %#v", PluginKey, pluginConfig)
	}
	return NewStore(cxt, cfg), nil
}
