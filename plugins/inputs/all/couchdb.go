//go:build !custom || inputs || inputs.couchdb

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/couchdb" // register plugin
