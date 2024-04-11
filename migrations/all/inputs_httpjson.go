//go:build !custom || (migrations && (inputs || inputs.httpjson))

package all

import _ "github.com/extremenetworks/telegraf/migrations/inputs_httpjson" // register migration
