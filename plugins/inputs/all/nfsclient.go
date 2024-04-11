//go:build !custom || inputs || inputs.nfsclient

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/nfsclient" // register plugin
