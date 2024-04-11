//go:build !custom || outputs || outputs.kafka

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/kafka" // register plugin
