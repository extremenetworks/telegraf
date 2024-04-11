//go:build !custom || inputs || inputs.sensors

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/sensors" // register plugin
