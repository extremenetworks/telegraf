//go:build !custom || inputs || inputs.ipset

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/ipset" // register plugin
