//go:build !custom || inputs || inputs.redis

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/redis" // register plugin
