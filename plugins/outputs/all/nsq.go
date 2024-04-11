//go:build !custom || outputs || outputs.nsq

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/nsq" // register plugin
