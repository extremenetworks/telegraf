//go:build !custom || inputs || inputs.nginx_plus_api

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/nginx_plus_api" // register plugin
