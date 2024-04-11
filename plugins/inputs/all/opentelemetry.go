//go:build !custom || inputs || inputs.opentelemetry

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/opentelemetry" // register plugin
