//go:build !custom || inputs || inputs.kapacitor

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/kapacitor" // register plugin
