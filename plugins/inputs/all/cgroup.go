//go:build !custom || inputs || inputs.cgroup

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/cgroup" // register plugin
