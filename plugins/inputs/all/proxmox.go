//go:build !custom || inputs || inputs.proxmox

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/proxmox" // register plugin
