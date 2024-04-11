//go:build !custom || inputs || inputs.netstat

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/netstat" // register plugin
