//go:build !custom || (migrations && (inputs || inputs.io || inputs.diskio))

package all

import _ "github.com/extremenetworks/telegraf/migrations/inputs_io" // register migration
