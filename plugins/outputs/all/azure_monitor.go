//go:build !custom || outputs || outputs.azure_monitor

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/azure_monitor" // register plugin
