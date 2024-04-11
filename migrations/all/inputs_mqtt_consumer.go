//go:build !custom || (migrations && (inputs || inputs.mqtt_consumer))

package all

import _ "github.com/extremenetworks/telegraf/migrations/inputs_mqtt_consumer" // register migration
