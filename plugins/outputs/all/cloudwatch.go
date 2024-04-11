//go:build !custom || outputs || outputs.cloudwatch

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/cloudwatch" // register plugin
