//go:build !custom || inputs || inputs.modbus

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/modbus" // register plugin
