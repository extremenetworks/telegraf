//go:build !custom || inputs || inputs.zipkin

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/zipkin" // register plugin
