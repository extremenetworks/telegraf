//go:build !custom || outputs || outputs.influxdb_v2

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/influxdb_v2" // register plugin
