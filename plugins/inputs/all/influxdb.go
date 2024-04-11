//go:build !custom || inputs || inputs.influxdb

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/influxdb" // register plugin
