//go:build !custom || outputs || outputs.nats

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/nats" // register plugin
