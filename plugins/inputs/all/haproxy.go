//go:build !custom || inputs || inputs.haproxy

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/haproxy" // register plugin
