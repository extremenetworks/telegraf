//go:build !custom || aggregators || aggregators.histogram

package all

import _ "github.com/extremenetworks/telegraf/plugins/aggregators/histogram" // register plugin
