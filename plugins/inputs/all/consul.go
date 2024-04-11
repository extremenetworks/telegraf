//go:build !custom || inputs || inputs.consul

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/consul" // register plugin
