//go:build !custom || inputs || inputs.mysql

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/mysql" // register plugin
