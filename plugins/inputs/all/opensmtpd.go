//go:build !custom || inputs || inputs.opensmtpd

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/opensmtpd" // register plugin
