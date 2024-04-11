//go:build !custom || inputs || inputs.ipmi_sensor

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/ipmi_sensor" // register plugin
