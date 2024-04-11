//go:build !custom || inputs || inputs.socketstat

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/socketstat" // register plugin
