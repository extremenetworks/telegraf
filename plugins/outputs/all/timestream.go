//go:build !custom || outputs || outputs.timestream

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/timestream" // register plugin
