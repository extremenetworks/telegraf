//go:build !custom || secretstores || secretstores.jose

package all

import _ "github.com/extremenetworks/telegraf/plugins/secretstores/jose" // register plugin
