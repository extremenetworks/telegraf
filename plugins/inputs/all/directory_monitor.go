//go:build !custom || inputs || inputs.directory_monitor

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/directory_monitor" // register plugin
