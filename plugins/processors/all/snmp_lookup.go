//go:build !custom || processors || processors.snmp_lookup

package all

import _ "github.com/extremenetworks/telegraf/plugins/processors/snmp_lookup" // register plugin
