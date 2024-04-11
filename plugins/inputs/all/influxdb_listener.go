//go:build !custom || inputs || inputs.influxdb_listener

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/influxdb_listener" // register plugin
