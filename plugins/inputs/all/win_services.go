//go:build !custom || inputs || inputs.win_services

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/win_services" // register plugin
