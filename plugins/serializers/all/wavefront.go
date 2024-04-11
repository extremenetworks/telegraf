//go:build !custom || serializers || serializers.wavefront

package all

import (
	_ "github.com/extremenetworks/telegraf/plugins/serializers/wavefront" // register plugin
)
