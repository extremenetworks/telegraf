//go:build !custom || outputs || outputs.cloud_pubsub

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/cloud_pubsub" // register plugin
