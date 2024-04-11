//go:build !custom || inputs || inputs.logparser

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/logparser" // register plugin
