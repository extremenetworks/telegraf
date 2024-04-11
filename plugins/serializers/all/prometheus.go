//go:build !custom || serializers || serializers.prometheus

package all

import (
	_ "github.com/extremenetworks/telegraf/plugins/serializers/prometheus" // register plugin
)
