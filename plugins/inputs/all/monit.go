//go:build !custom || inputs || inputs.monit

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/monit" // register plugin
