//go:build !custom || inputs || inputs.intel_pmu

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/intel_pmu" // register plugin
