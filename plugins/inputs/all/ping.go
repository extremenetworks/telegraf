//go:build !custom || inputs || inputs.ping

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/ping" // register plugin
