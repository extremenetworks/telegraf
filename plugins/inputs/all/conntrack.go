//go:build !custom || inputs || inputs.conntrack

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/conntrack" // register plugin
