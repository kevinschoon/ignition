package linuxkit

import (
	"github.com/coreos/ignition/config/types"
	"github.com/coreos/ignition/internal/log"
	"github.com/coreos/ignition/internal/providers"
	"github.com/coreos/ignition/internal/resource"
)

type LinuxKit struct {
	name              string
	fetch             providers.FuncFetchConfig
	newFetcher        providers.FuncNewFetcher
	baseConfig        types.Config
	defaultUserConfig types.Config
}

func (c LinuxKit) Name() string {
	return c.name
}

func (c LinuxKit) FetchFunc() providers.FuncFetchConfig {
	return c.fetch
}

func (c LinuxKit) NewFetcherFunc() providers.FuncNewFetcher {
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

func (c LinuxKit) BaseConfig() types.Config {
	return c.baseConfig
}

func (c LinuxKit) DefaultUserConfig() types.Config {
	return c.defaultUserConfig
}
