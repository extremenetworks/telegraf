//go:build !custom || inputs || inputs.win_eventlog

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/win_eventlog" // register plugin
