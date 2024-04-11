//go:build !custom || serializers || serializers.json

package all

import (
	_ "github.com/extremenetworks/telegraf/plugins/serializers/json" // register plugin
)
