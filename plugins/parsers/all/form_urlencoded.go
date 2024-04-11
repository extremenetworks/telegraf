//go:build !custom || parsers || parsers.form_urlencoded

package all

import _ "github.com/extremenetworks/telegraf/plugins/parsers/form_urlencoded" // register plugin
