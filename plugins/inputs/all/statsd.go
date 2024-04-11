//go:build !custom || inputs || inputs.statsd

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/statsd" // register plugin
