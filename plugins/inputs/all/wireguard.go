//go:build !custom || inputs || inputs.wireguard

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/wireguard" // register plugin
