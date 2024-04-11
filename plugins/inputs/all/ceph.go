//go:build !custom || inputs || inputs.ceph

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/ceph" // register plugin
