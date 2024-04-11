//go:build !custom || parsers || parsers.prometheus

package all

import _ "github.com/extremenetworks/telegraf/plugins/parsers/prometheus" // register plugin
