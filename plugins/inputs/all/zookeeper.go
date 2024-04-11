//go:build !custom || inputs || inputs.zookeeper

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/zookeeper" // register plugin
