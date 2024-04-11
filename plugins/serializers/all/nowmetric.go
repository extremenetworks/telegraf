//go:build !custom || serializers || serializers.nowmetric

package all

import (
	_ "github.com/extremenetworks/telegraf/plugins/serializers/nowmetric" // register plugin
)
