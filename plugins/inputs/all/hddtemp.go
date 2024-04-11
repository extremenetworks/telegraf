//go:build !custom || inputs || inputs.hddtemp

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/hddtemp" // register plugin
