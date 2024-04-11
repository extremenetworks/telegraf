//go:build !custom || inputs || inputs.apcupsd

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/apcupsd" // register plugin
