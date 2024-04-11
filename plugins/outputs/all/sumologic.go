//go:build !custom || outputs || outputs.sumologic

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/sumologic" // register plugin
