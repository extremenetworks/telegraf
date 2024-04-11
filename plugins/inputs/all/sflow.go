//go:build !custom || inputs || inputs.sflow

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/sflow" // register plugin
