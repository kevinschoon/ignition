package linuxkit

import (
	//"github.com/coreos/ignition/config/types"
	"github.com/coreos/ignition/internal/providers/azure"
	"github.com/coreos/ignition/internal/providers/digitalocean"
	"github.com/coreos/ignition/internal/providers/ec2"
	"github.com/coreos/ignition/internal/providers/file"
	"github.com/coreos/ignition/internal/providers/gce"
	"github.com/coreos/ignition/internal/providers/noop"
	"github.com/coreos/ignition/internal/providers/openstack"
	"github.com/coreos/ignition/internal/providers/packet"
	"github.com/coreos/ignition/internal/providers/qemu"
	"github.com/coreos/ignition/internal/providers/virtualbox"
	"github.com/coreos/ignition/internal/providers/vmware"
	"github.com/coreos/ignition/internal/registry"
	"github.com/coreos/ignition/internal/util"
)

var Config = registry.Create("oem configs")
var yes = util.BoolToPtr(true)

func init() {
	Config.Register(LinuxKit{
		name:  "azure",
		fetch: azure.FetchConfig,
	})
	Config.Register(LinuxKit{
		name:  "cloudsigma",
		fetch: noop.FetchConfig,
	})
	Config.Register(LinuxKit{
		name:  "cloudstack",
		fetch: noop.FetchConfig,
	})
	Config.Register(LinuxKit{
		name:  "digitalocean",
		fetch: digitalocean.FetchConfig,
	})
	Config.Register(LinuxKit{
		name:  "brightbox",
		fetch: openstack.FetchConfig,
	})
	Config.Register(LinuxKit{
		name:  "openstack",
		fetch: openstack.FetchConfig,
	})
	Config.Register(LinuxKit{
		name:       "ec2",
		fetch:      ec2.FetchConfig,
		newFetcher: ec2.NewFetcher,
	})
	Config.Register(LinuxKit{
		name:  "exoscale",
		fetch: noop.FetchConfig,
	})
	Config.Register(LinuxKit{
		name:  "gce",
		fetch: gce.FetchConfig,
	})
	Config.Register(LinuxKit{
		name:  "hyperv",
		fetch: noop.FetchConfig,
	})
	Config.Register(LinuxKit{
		name:  "niftycloud",
		fetch: noop.FetchConfig,
	})
	Config.Register(LinuxKit{
		name:  "packet",
		fetch: packet.FetchConfig,
	})
	Config.Register(LinuxKit{
		name:  "pxe",
		fetch: noop.FetchConfig,
	})
	Config.Register(LinuxKit{
		name:  "rackspace",
		fetch: noop.FetchConfig,
	})
	Config.Register(LinuxKit{
		name:  "rackspace-onmetal",
		fetch: noop.FetchConfig,
	})
	Config.Register(LinuxKit{
		name:  "vagrant",
		fetch: noop.FetchConfig,
	})
	Config.Register(LinuxKit{
		name:  "vagrant-virtualbox",
		fetch: virtualbox.FetchConfig,
	})
	Config.Register(LinuxKit{
		name:  "virtualbox",
		fetch: virtualbox.FetchConfig,
	})
	Config.Register(LinuxKit{
		name:  "vmware",
		fetch: vmware.FetchConfig,
	})
	Config.Register(LinuxKit{
		name:  "interoute",
		fetch: noop.FetchConfig,
	})
	Config.Register(LinuxKit{
		name:  "qemu",
		fetch: qemu.FetchConfig,
	})
	Config.Register(LinuxKit{
		name:  "file",
		fetch: file.FetchConfig,
	})
}
