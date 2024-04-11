//go:build !custom || inputs || inputs.kube_inventory

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/kube_inventory" // register plugin
