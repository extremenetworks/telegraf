//go:build !custom || inputs || inputs.win_perf_counters

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/win_perf_counters" // register plugin
