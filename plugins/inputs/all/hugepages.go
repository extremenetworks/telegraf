//go:build !custom || inputs || inputs.hugepages

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/hugepages" // register plugin
