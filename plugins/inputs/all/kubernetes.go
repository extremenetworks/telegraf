//go:build !custom || inputs || inputs.kubernetes

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/kubernetes" // register plugin
