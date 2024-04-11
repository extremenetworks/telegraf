//go:build !custom || inputs || inputs.nvidia_smi

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/nvidia_smi" // register plugin
