//go:build !custom || (migrations && (inputs || inputs.disk))

package all

import _ "github.com/extremenetworks/telegraf/migrations/inputs_disk" // register migration
