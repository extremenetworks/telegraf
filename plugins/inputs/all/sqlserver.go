//go:build !custom || inputs || inputs.sqlserver

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/sqlserver" // register plugin
