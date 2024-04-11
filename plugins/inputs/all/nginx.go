//go:build !custom || inputs || inputs.nginx

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/nginx" // register plugin
