//go:build !custom || inputs || inputs.cloud_pubsub

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/cloud_pubsub" // register plugin
