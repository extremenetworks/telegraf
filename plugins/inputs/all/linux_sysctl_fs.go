//go:build !custom || inputs || inputs.linux_sysctl_fs

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/linux_sysctl_fs" // register plugin
