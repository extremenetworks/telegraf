//go:build !custom || inputs || inputs.docker_log

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/docker_log" // register plugin
