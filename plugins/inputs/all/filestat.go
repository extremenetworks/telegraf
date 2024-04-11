//go:build !custom || inputs || inputs.filestat

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/filestat" // register plugin
