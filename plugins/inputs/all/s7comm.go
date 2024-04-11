//go:build !custom || inputs || inputs.s7comm

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/s7comm" // register plugin
