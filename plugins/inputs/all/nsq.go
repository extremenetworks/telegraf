//go:build !custom || inputs || inputs.nsq

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/nsq" // register plugin
