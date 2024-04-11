//go:build !custom || inputs || inputs.phpfpm

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/phpfpm" // register plugin
