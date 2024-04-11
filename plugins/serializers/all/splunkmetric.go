//go:build !custom || serializers || serializers.splunkmetric

package all

import (
	_ "github.com/extremenetworks/telegraf/plugins/serializers/splunkmetric" // register plugin
)
