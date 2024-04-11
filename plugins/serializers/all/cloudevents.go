//go:build !custom || serializers || serializers.cloudevents

package all

import (
	_ "github.com/extremenetworks/telegraf/plugins/serializers/cloudevents" // register plugin
)
