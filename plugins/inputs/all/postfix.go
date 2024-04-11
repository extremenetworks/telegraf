//go:build !custom || inputs || inputs.postfix

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/postfix" // register plugin
