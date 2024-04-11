//go:build !custom || inputs || inputs.docker

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/docker" // register plugin
