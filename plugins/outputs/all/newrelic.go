//go:build !custom || outputs || outputs.newrelic

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/newrelic" // register plugin
