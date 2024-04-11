//go:build !custom || inputs || inputs.webhooks

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/webhooks" // register plugin
