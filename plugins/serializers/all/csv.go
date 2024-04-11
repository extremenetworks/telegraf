//go:build !custom || serializers || serializers.csv

package all

import (
	_ "github.com/extremenetworks/telegraf/plugins/serializers/csv" // register plugin
)
