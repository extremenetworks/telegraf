//go:build !custom || outputs || outputs.bigquery

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/bigquery" // register plugin
