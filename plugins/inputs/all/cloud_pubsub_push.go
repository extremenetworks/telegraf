//go:build !custom || inputs || inputs.cloud_pubsub_push

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/cloud_pubsub_push" // register plugin
