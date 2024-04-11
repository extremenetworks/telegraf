//go:build !custom || inputs || inputs.nats

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/nats" // register plugin
