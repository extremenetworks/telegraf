//go:build !custom || inputs || inputs.win_wmi

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/win_wmi" // register plugin
