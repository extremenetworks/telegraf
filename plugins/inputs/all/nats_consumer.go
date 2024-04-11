//go:build !custom || inputs || inputs.nats_consumer

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/nats_consumer" // register plugin
