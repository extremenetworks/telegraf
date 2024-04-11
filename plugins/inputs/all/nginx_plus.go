//go:build !custom || inputs || inputs.nginx_plus

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/nginx_plus" // register plugin
