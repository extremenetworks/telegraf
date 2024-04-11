//go:build !custom || outputs || outputs.riemann

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/riemann" // register plugin
