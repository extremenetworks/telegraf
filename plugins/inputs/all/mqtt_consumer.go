//go:build !custom || inputs || inputs.mqtt_consumer

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/mqtt_consumer" // register plugin
