//go:build !custom || outputs || outputs.zabbix

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/zabbix" // register plugin
