//go:build !custom || inputs || inputs.prometheus

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/prometheus" // register plugin
