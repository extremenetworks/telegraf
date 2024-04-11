//go:build !custom || inputs || inputs.rethinkdb

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/rethinkdb" // register plugin
