//go:build !custom || processors || processors.lookup

package all

import _ "github.com/extremenetworks/telegraf/plugins/processors/lookup" // register plugin
