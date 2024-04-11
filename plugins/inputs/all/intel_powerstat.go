//go:build !custom || inputs || inputs.intel_powerstat

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/intel_powerstat" // register plugin
