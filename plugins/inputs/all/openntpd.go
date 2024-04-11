//go:build !custom || inputs || inputs.openntpd

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/openntpd" // register plugin
