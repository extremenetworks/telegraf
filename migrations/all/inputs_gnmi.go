//go:build !custom || (migrations && (inputs || inputs.gnmi))

package all

import _ "github.com/extremenetworks/telegraf/migrations/inputs_gnmi" // register migration
