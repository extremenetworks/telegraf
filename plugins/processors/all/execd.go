//go:build !custom || processors || processors.execd

package all

import _ "github.com/extremenetworks/telegraf/plugins/processors/execd" // register plugin
