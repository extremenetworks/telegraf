//go:build !custom || inputs || inputs.dns_query

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/dns_query" // register plugin
