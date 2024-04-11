//go:build !custom || outputs || outputs.datadog

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/datadog" // register plugin
