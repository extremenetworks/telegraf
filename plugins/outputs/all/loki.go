//go:build !custom || outputs || outputs.loki

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/loki" // register plugin
