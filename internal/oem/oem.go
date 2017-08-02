// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package oem

import (
	"fmt"

	"github.com/coreos/ignition/config/types"
	"github.com/coreos/ignition/internal/oem/coreos"
	"github.com/coreos/ignition/internal/providers"
)

type Options struct {
	Name     string
	Platform string
}

func (o Options) String() string { return fmt.Sprintf("%s %s", o.Name, o.Platform) }

// Config represents a set of options that map to a particular OEM and platform
type Config interface {
	Name() string
	FetchFunc() providers.FuncFetchConfig
	NewFetcherFunc() providers.FuncNewFetcher
	BaseConfig() types.Config
	DefaultUserConfig() types.Config
}

func Get(opts Options) (config Config, ok bool) {
	switch opts.Platform {
	case "coreos":
		config, ok = coreos.Config.Get(opts.Name).(Config)
		return
	}
	return
}

func MustGet(opts Options) Config {
	if config, ok := Get(opts); ok {
		return config
	} else {
		panic(fmt.Sprintf("invalid Platform/OEM %s %s provided", opts.Name, opts.Platform))
	}
}

func Names(opts Options) (names []string) {
	switch opts.Platform {
	case "coreos":
		names = coreos.Config.Names()
	}
	return
}
