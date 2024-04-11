//go:build !custom || processors || processors.tag_limit

package all

import _ "github.com/extremenetworks/telegraf/plugins/processors/tag_limit" // register plugin
