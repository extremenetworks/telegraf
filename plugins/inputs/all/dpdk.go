//go:build !custom || inputs || inputs.dpdk

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/dpdk" // register plugin
