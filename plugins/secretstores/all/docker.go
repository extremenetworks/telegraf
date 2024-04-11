//go:build !custom || secretstores || secretstores.docker

package all

import _ "github.com/extremenetworks/telegraf/plugins/secretstores/docker" // register plugin
