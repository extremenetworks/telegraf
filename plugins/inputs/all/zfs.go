//go:build !custom || inputs || inputs.zfs

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/zfs" // register plugin
