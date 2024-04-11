//go:build !custom || outputs || outputs.syslog

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/syslog" // register plugin
