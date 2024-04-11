//go:build !custom || inputs || inputs.cloudwatch

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/cloudwatch" // register plugin
