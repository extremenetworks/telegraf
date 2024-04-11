//go:build !custom || inputs || inputs.chrony

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/chrony" // register plugin
