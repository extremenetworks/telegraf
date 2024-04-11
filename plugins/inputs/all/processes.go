//go:build !custom || inputs || inputs.processes

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/processes" // register plugin
