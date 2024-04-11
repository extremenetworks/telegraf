//go:build !custom || outputs || outputs.azure_data_explorer

package all

import _ "github.com/extremenetworks/telegraf/plugins/outputs/azure_data_explorer" // register plugin
