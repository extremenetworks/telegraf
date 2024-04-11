//go:build !custom || parsers || parsers.nagios

package all

import _ "github.com/extremenetworks/telegraf/plugins/parsers/nagios" // register plugin
