//go:build !custom || outputs || outputs.postgresql

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/postgresql" // register plugin
