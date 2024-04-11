//go:build !custom || outputs || outputs.event_hubs

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/event_hubs" // register plugin
