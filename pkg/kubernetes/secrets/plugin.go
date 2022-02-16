package secrets

import (
	portercontext "get.porter.sh/porter/pkg/context"
	"get.porter.sh/porter/pkg/secrets"
	"get.porter.sh/porter/pkg/secrets/plugins"
	cnabsecrets "github.com/cnabio/cnab-go/secrets"
	hplugin "github.com/hashicorp/go-plugin"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

const PluginKey = plugins.PluginInterface + ".kubernetes.secret"

var _ cnabsecrets.Store = &Plugin{}

type PluginConfig struct {
	KubeConfig string `mapstructure:"kubeconfig"`
	Namespace  string `mapstructure:"namespace"`
}

type Plugin struct {
	cnabsecrets.Store
}

func NewPlugin(cxt *portercontext.Context, pluginConfig interface{}) (hplugin.Plugin, error) {
	cfg := PluginConfig{}

	if err := mapstructure.Decode(pluginConfig, &cfg); err != nil {
		return nil, errors.Wrapf(err, "error decoding %s plugin config from %#v", PluginKey, pluginConfig)
	}
	return &secrets.Plugin{
		Impl: &Plugin{
			Store: NewStore(cxt, cfg),
		},
	}, nil
}
