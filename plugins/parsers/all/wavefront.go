//go:build !custom || parsers || parsers.wavefront

package all

import _ "github.com/extremenetworks/telegraf/plugins/parsers/wavefront" // register plugin
