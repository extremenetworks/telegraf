//go:build !custom || outputs || outputs.amqp

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/amqp" // register plugin
