//go:build !custom || outputs || outputs.opentelemetry

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/opentelemetry" // register plugin
