//go:build !custom || secretstores || secretstores.systemd

package all

import _ "github.com/extremenetworks/telegraf/plugins/secretstores/systemd" // register plugin
