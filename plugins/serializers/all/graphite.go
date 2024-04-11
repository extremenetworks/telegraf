//go:build !custom || serializers || serializers.graphite

package all

import (
	_ "github.com/extremenetworks/telegraf/plugins/serializers/graphite" // register plugin
)
