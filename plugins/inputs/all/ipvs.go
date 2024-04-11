//go:build !custom || inputs || inputs.ipvs

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/ipvs" // register plugin
