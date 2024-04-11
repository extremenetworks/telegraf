//go:build !custom || parsers || parsers.csv

package all

import _ "github.com/extremenetworks/telegraf/plugins/parsers/csv" // register plugin
