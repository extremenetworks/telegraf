//go:build !custom || inputs || inputs.nomad

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/nomad" // register plugin
