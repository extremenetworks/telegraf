//go:build !custom || processors || processors.reverse_dns

package all

import _ "github.com/extremenetworks/telegraf/plugins/processors/reverse_dns" // register plugin
