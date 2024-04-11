//go:build !custom || inputs || inputs.ldap

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/ldap" // register plugin
