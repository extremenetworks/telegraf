//go:build !custom || secretstores || secretstores.oauth2

package all

import _ "github.com/extremenetworks/telegraf/plugins/secretstores/oauth2" // register plugin
