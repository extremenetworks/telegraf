//go:build !custom || parsers || parsers.opentsdb

package all

import _ "github.com/extremenetworks/telegraf/plugins/parsers/opentsdb" // register plugin
