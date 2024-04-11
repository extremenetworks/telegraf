//go:build !custom || outputs || outputs.stackdriver

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/stackdriver" // register plugin
