//go:build !custom || outputs || outputs.wavefront

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/wavefront" // register plugin
