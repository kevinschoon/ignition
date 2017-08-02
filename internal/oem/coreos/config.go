package coreos

import (
	"github.com/coreos/ignition/config/types"
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

	"fmt"
	"github.com/vincent-petithory/dataurl"
	"net/url"
)

var Config = registry.Create("oem configs")
var yes = util.BoolToPtr(true)

func init() {
	Config.Register(CoreOS{
		name:  "azure",
		fetch: azure.FetchConfig,
		baseConfig: types.Config{
			Systemd: types.Systemd{
				Units: []types.Unit{
					{Enabled: yes, Name: "waagent.service"},
					{Name: "etcd2.service", Dropins: []types.Dropin{
						{Name: "10-oem.conf", Contents: "[Service]\nEnvironment=ETCD_ELECTION_TIMEOUT=1200\n"},
					}},
				},
			},
			Storage: types.Storage{Files: []types.File{serviceFromOem("waagent.service")}},
		},
		defaultUserConfig: types.Config{Systemd: types.Systemd{Units: []types.Unit{userCloudInit("Azure", "azure")}}},
	})
	Config.Register(CoreOS{
		name:  "cloudsigma",
		fetch: noop.FetchConfig,
	})
	Config.Register(CoreOS{
		name:  "cloudstack",
		fetch: noop.FetchConfig,
	})
	Config.Register(CoreOS{
		name:  "digitalocean",
		fetch: digitalocean.FetchConfig,
		baseConfig: types.Config{
			Systemd: types.Systemd{
				Units: []types.Unit{{Enabled: yes, Name: "coreos-metadata-sshkeys@.service"}},
			},
		},
		defaultUserConfig: types.Config{Systemd: types.Systemd{Units: []types.Unit{userCloudInit("DigitalOcean", "digitalocean")}}},
	})
	Config.Register(CoreOS{
		name:  "brightbox",
		fetch: openstack.FetchConfig,
		defaultUserConfig: types.Config{
			Systemd: types.Systemd{
				Units: []types.Unit{
					{Mask: true, Name: "user-configdrive.service"},
					{Mask: true, Name: "user-configvirtfs.service"},
					userCloudInit("BrightBox", "ec2-compat"),
				},
			},
		},
	})
	Config.Register(CoreOS{
		name:  "openstack",
		fetch: openstack.FetchConfig,
		defaultUserConfig: types.Config{
			Systemd: types.Systemd{
				Units: []types.Unit{
					{Mask: true, Name: "user-configdrive.service"},
					{Mask: true, Name: "user-configvirtfs.service"},
					userCloudInit("OpenStack", "ec2-compat"),
				},
			},
		},
	})
	Config.Register(CoreOS{
		name:       "ec2",
		fetch:      ec2.FetchConfig,
		newFetcher: ec2.NewFetcher,
		baseConfig: types.Config{
			Systemd: types.Systemd{
				Units: []types.Unit{
					{Enabled: yes, Name: "coreos-metadata-sshkeys@.service"},
					{Name: "etcd2.service", Dropins: []types.Dropin{
						{Name: "10-oem.conf", Contents: "[Service]\nEnvironment=ETCD_ELECTION_TIMEOUT=1200\n"},
					}},
				},
			},
		},
		defaultUserConfig: types.Config{
			Systemd: types.Systemd{
				Units: []types.Unit{
					{Mask: true, Name: "user-configdrive.service"},
					{Mask: true, Name: "user-configvirtfs.service"},
					userCloudInit("EC2", "ec2-compat"),
				},
			},
		},
	})
	Config.Register(CoreOS{
		name:  "exoscale",
		fetch: noop.FetchConfig,
	})
	Config.Register(CoreOS{
		name:  "gce",
		fetch: gce.FetchConfig,
		baseConfig: types.Config{
			Systemd: types.Systemd{
				Units: []types.Unit{
					{Enabled: yes, Name: "coreos-metadata-sshkeys@.service"},
					{Enabled: yes, Name: "oem-gce.service"},
				},
			},
			Storage: types.Storage{
				Files: []types.File{
					serviceFromOem("oem-gce.service"),
					{
						Node: types.Node{
							Filesystem: "root",
							Path:       "/etc/hosts",
						},
						FileEmbedded1: types.FileEmbedded1{
							Mode:     0444,
							Contents: contentsFromString("169.254.169.254 metadata\n127.0.0.1 localhost\n"),
						},
					},
					{
						Node: types.Node{
							Filesystem: "root",
							Path:       "/etc/profile.d/google-cloud-sdk.sh",
						},
						FileEmbedded1: types.FileEmbedded1{
							Mode: 0444,
							Contents: contentsFromString(`#!/bin/sh
alias gcloud="(docker images google/cloud-sdk || docker pull google/cloud-sdk) > /dev/null;docker run -t -i --net="host" -v $HOME/.config:/.config -v /var/run/docker.sock:/var/run/doker.sock -v /usr/bin/docker:/usr/bin/docker google/cloud-sdk gcloud"
alias gcutil="(docker images google/cloud-sdk || docker pull google/cloud-sdk) > /dev/null;docker run -t -i --net="host" -v $HOME/.config:/.config google/cloud-sdk gcutil"
alias gsutil="(docker images google/cloud-sdk || docker pull google/cloud-sdk) > /dev/null;docker run -t -i --net="host" -v $HOME/.config:/.config google/cloud-sdk gsutil"
`),
						},
					},
				},
			},
		},
		defaultUserConfig: types.Config{Systemd: types.Systemd{Units: []types.Unit{userCloudInit("GCE", "gce")}}},
	})
	Config.Register(CoreOS{
		name:  "hyperv",
		fetch: noop.FetchConfig,
	})
	Config.Register(CoreOS{
		name:  "niftycloud",
		fetch: noop.FetchConfig,
	})
	Config.Register(CoreOS{
		name:  "packet",
		fetch: packet.FetchConfig,
		baseConfig: types.Config{
			Systemd: types.Systemd{
				Units: []types.Unit{
					{Enabled: yes, Name: "coreos-metadata-sshkeys@.service"},
					{Enabled: yes, Name: "packet-phone-home.service"},
				},
			},
			Storage: types.Storage{Files: []types.File{serviceFromOem("packet-phone-home.service")}},
		},
		defaultUserConfig: types.Config{Systemd: types.Systemd{Units: []types.Unit{userCloudInit("Packet", "packet")}}},
	})
	Config.Register(CoreOS{
		name:  "pxe",
		fetch: noop.FetchConfig,
	})
	Config.Register(CoreOS{
		name:  "rackspace",
		fetch: noop.FetchConfig,
	})
	Config.Register(CoreOS{
		name:  "rackspace-onmetal",
		fetch: noop.FetchConfig,
	})
	Config.Register(CoreOS{
		name:  "vagrant",
		fetch: noop.FetchConfig,
	})
	Config.Register(CoreOS{
		name:  "vagrant-virtualbox",
		fetch: virtualbox.FetchConfig,
	})
	Config.Register(CoreOS{
		name:  "virtualbox",
		fetch: virtualbox.FetchConfig,
	})
	Config.Register(CoreOS{
		name:  "vmware",
		fetch: vmware.FetchConfig,
		baseConfig: types.Config{
			Systemd: types.Systemd{Units: []types.Unit{{Enabled: yes, Name: "vmtoolsd.service"}}},
			Storage: types.Storage{Files: []types.File{serviceFromOem("vmtoolsd.service")}},
		},
		defaultUserConfig: types.Config{Systemd: types.Systemd{Units: []types.Unit{userCloudInit("VMware", "vmware")}}},
	})
	Config.Register(CoreOS{
		name:  "interoute",
		fetch: noop.FetchConfig,
	})
	Config.Register(CoreOS{
		name:  "qemu",
		fetch: qemu.FetchConfig,
	})
	Config.Register(CoreOS{
		name:  "file",
		fetch: file.FetchConfig,
	})
}

func contentsFromString(data string) types.FileContents {
	return types.FileContents{
		Source: (&url.URL{
			Scheme: "data",
			Opaque: "," + dataurl.EscapeString(data),
		}).String(),
	}
}

func contentsFromOem(path string) types.FileContents {
	return types.FileContents{
		Source: (&url.URL{
			Scheme: "oem",
			Path:   path,
		}).String(),
	}
}

func userCloudInit(name string, oem string) types.Unit {
	contents := `[Unit]
Description=Cloudinit from %s metadata

[Service]
Type=oneshot
ExecStart=/usr/bin/coreos-cloudinit --oem=%s

[Install]
WantedBy=multi-user.target
`

	return types.Unit{
		Name:     "oem-cloudinit.service",
		Enabled:  yes,
		Contents: fmt.Sprintf(contents, name, oem),
	}
}

func serviceFromOem(unit string) types.File {
	return types.File{
		Node: types.Node{
			Filesystem: "root",
			Path:       "/etc/systemd/system/" + unit,
		},
		FileEmbedded1: types.FileEmbedded1{
			Mode:     0444,
			Contents: contentsFromOem("/units/" + unit),
		},
	}
}
