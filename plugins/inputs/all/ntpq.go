//go:build !custom || inputs || inputs.ntpq

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/ntpq" // register plugin
