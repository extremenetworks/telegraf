//go:build !custom || parsers || parsers.prometheusremotewrite

package all

import _ "github.com/extremenetworks/telegraf/plugins/parsers/prometheusremotewrite" // register plugin
