//go:build !custom || inputs || inputs.postgresql_extensible

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/postgresql_extensible" // register plugin
