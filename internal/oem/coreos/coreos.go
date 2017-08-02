package coreos

import (
	"github.com/coreos/ignition/config/types"
	"github.com/coreos/ignition/internal/log"
	"github.com/coreos/ignition/internal/providers"
	"github.com/coreos/ignition/internal/resource"
)

type CoreOS struct {
	name              string
	fetch             providers.FuncFetchConfig
	newFetcher        providers.FuncNewFetcher
	baseConfig        types.Config
	defaultUserConfig types.Config
}

func (c CoreOS) Name() string {
	return c.name
}

func (c CoreOS) FetchFunc() providers.FuncFetchConfig {
	return c.fetch
}

func (c CoreOS) NewFetcherFunc() providers.FuncNewFetcher {
	if c.newFetcher != nil {
		return c.newFetcher
	}
	return func(l *log.Logger, c *resource.HttpClient) (resource.Fetcher, error) {
		return resource.Fetcher{
			Logger: l,
			Client: c,
		}, nil
	}
}

func (c CoreOS) BaseConfig() types.Config {
	return c.baseConfig
}

func (c CoreOS) DefaultUserConfig() types.Config {
	return c.defaultUserConfig
}
