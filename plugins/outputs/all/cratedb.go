//go:build !custom || outputs || outputs.cratedb

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/cratedb" // register plugin
