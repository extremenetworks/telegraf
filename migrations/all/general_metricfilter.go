//go:build !custom || migrations

package all

import _ "github.com/extremenetworks/telegraf/migrations/general_metricfilter" // register migration
