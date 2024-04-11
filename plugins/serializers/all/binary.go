//go:build !custom || serializers || serializers.binary

package all

import (
	_ "github.com/extremenetworks/telegraf/plugins/serializers/binary" // register plugin
)
