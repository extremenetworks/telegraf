//go:build !custom || inputs || inputs.kinesis_consumer

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/kinesis_consumer" // register plugin
