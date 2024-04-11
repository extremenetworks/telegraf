//go:build !custom || inputs || inputs.ethtool

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/ethtool" // register plugin
