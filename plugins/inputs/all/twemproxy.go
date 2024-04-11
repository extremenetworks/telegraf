//go:build !custom || inputs || inputs.twemproxy

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/twemproxy" // register plugin
