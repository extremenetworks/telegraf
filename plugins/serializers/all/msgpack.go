//go:build !custom || serializers || serializers.msgpack

package all

import (
	_ "github.com/extremenetworks/telegraf/plugins/serializers/msgpack" // register plugin
)
