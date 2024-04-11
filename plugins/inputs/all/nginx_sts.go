//go:build !custom || inputs || inputs.nginx_sts

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/nginx_sts" // register plugin
