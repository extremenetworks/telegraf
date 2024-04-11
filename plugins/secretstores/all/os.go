//go:build !custom || secretstores || secretstores.os

package all

import _ "github.com/extremenetworks/telegraf/plugins/secretstores/os" // register plugin
