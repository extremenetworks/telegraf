//go:build !custom || serializers || serializers.influx

package all

import (
	_ "github.com/extremenetworks/telegraf/plugins/serializers/influx" // register plugin
)
