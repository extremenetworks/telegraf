//go:build !custom || inputs || inputs.kernel_vmstat

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/kernel_vmstat" // register plugin
