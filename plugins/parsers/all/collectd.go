//go:build !custom || parsers || parsers.collectd

package all

import _ "github.com/extremenetworks/telegraf/plugins/parsers/collectd" // register plugin
