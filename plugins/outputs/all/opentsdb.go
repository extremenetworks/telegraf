//go:build !custom || outputs || outputs.opentsdb

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/opentsdb" // register plugin
