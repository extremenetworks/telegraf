//go:build !custom || inputs || inputs.snmp

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/snmp" // register plugin
