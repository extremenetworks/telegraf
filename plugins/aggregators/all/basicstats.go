//go:build !custom || aggregators || aggregators.basicstats

package all

import _ "github.com/extremenetworks/telegraf/plugins/aggregators/basicstats" // register plugin
