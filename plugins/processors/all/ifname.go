//go:build !custom || processors || processors.ifname

package all

import _ "github.com/extremenetworks/telegraf/plugins/processors/ifname" // register plugin
