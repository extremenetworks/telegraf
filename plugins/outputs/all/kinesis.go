//go:build !custom || outputs || outputs.kinesis

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/kinesis" // register plugin
