//go:build !custom || outputs || outputs.sql

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/sql" // register plugin
