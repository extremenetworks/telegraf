//go:build !custom || inputs || inputs.openweathermap

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/openweathermap" // register plugin
