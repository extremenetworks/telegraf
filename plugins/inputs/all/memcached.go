//go:build !custom || inputs || inputs.memcached

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/memcached" // register plugin
