//go:build !custom || processors || processors.dedup

package all

import _ "github.com/extremenetworks/telegraf/plugins/processors/dedup" // register plugin
