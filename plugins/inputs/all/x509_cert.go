//go:build !custom || inputs || inputs.x509_cert

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/x509_cert" // register plugin
