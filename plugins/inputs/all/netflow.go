//go:build !custom || inputs || inputs.netflow

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/netflow" // register plugin
