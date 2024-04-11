//go:build !custom || outputs || outputs.socket_writer

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/socket_writer" // register plugin
