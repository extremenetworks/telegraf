//go:build !custom || inputs || inputs.eventhub_consumer

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/eventhub_consumer" // register plugin
