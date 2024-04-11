//go:build !custom || inputs || inputs.http_listener_v2

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/http_listener_v2" // register plugin
