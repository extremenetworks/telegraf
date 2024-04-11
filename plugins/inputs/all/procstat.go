//go:build !custom || inputs || inputs.procstat

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/procstat" // register plugin
