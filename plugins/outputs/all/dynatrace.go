//go:build !custom || outputs || outputs.dynatrace

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/dynatrace" // register plugin
