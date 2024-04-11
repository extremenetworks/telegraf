//go:build !custom || inputs || inputs.azure_storage_queue

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/azure_storage_queue" // register plugin
