//go:build !custom || inputs || inputs.iptables

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/iptables" // register plugin
