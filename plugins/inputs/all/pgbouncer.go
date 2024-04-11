//go:build !custom || inputs || inputs.pgbouncer

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/pgbouncer" // register plugin
