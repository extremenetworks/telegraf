//go:build !custom || inputs || inputs.riemann_listener

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/riemann_listener" // register plugin
