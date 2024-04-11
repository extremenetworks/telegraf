//go:build !custom || inputs || inputs.consul_agent

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/consul_agent" // register plugin
