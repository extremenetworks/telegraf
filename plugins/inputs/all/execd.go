//go:build !custom || inputs || inputs.execd

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/execd" // register plugin
