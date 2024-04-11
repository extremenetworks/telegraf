//go:build !custom || inputs || inputs.clickhouse

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/clickhouse" // register plugin
