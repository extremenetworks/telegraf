//go:build !custom || parsers || parsers.logfmt

package all

import _ "github.com/extremenetworks/telegraf/plugins/parsers/logfmt" // register plugin
