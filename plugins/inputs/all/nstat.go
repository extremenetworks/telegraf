//go:build !custom || inputs || inputs.nstat

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/nstat" // register plugin
