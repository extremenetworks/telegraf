//go:build !custom || inputs || inputs.sysstat

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/sysstat" // register plugin
