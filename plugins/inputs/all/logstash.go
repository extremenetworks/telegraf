//go:build !custom || inputs || inputs.logstash

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/logstash" // register plugin
