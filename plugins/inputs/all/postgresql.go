//go:build !custom || inputs || inputs.postgresql

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/postgresql" // register plugin
