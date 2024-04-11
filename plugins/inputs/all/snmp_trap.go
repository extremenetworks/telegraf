//go:build !custom || inputs || inputs.snmp_trap

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/snmp_trap" // register plugin
