//go:build !custom || inputs || inputs.amd_rocm_smi

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/amd_rocm_smi" // register plugin
