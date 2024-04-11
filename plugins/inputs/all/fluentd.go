//go:build !custom || inputs || inputs.fluentd

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/fluentd" // register plugin
