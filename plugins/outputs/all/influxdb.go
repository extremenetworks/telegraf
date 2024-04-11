//go:build !custom || outputs || outputs.influxdb

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/influxdb" // register plugin
