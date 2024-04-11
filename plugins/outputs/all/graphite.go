//go:build !custom || outputs || outputs.graphite

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/graphite" // register plugin
