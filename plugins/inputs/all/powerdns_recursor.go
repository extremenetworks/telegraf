//go:build !custom || inputs || inputs.powerdns_recursor

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/powerdns_recursor" // register plugin
