//go:build !custom || outputs || outputs.sensu

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/sensu" // register plugin
