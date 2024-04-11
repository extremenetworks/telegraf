//go:build !custom || serializers || serializers.prometheusremotewrite

package all

import (
	_ "github.com/extremenetworks/telegraf/plugins/serializers/prometheusremotewrite" // register plugin
)
