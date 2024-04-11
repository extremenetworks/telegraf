package outputs

import "github.com/extremenetworks/telegraf"

// Deprecations lists the deprecated plugins
var Deprecations = map[string]telegraf.DeprecationInfo{
	"riemann_legacy": {
		Since:     "1.3.0",
		RemovalIn: "1.30.0",
		Notice:    "use 'outputs.riemann' instead (see https://github.com/extremenetworks/telegraf/issues/1878)",
	},
}
