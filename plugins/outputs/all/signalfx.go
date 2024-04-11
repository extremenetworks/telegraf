//go:build !custom || outputs || outputs.signalfx

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/signalfx" // register plugin
