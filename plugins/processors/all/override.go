//go:build !custom || processors || processors.override

package all

import _ "github.com/extremenetworks/telegraf/plugins/processors/override" // register plugin
