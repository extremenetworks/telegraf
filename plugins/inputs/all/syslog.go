//go:build !custom || inputs || inputs.syslog

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/syslog" // register plugin
