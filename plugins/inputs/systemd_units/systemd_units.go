//go:generate ../../../tools/readme_config_includer/generator
package systemd_units

import (
	_ "embed"
	"time"

	"github.com/extremenetworks/telegraf"
	"github.com/extremenetworks/telegraf/config"
	"github.com/extremenetworks/telegraf/plugins/inputs"
)

//go:embed sample.conf
var sampleConfig string

// SystemdUnits is a telegraf plugin to gather systemd unit status
type SystemdUnits struct {
	Pattern  string          `toml:"pattern"`
	UnitType string          `toml:"unittype"`
	Details  bool            `toml:"details"`
	Timeout  config.Duration `toml:"timeout"`
	Log      telegraf.Logger `toml:"-"`
	archParams
}

func (*SystemdUnits) SampleConfig() string {
	return sampleConfig
}

func init() {
	inputs.Add("systemd_units", func() telegraf.Input {
		return &SystemdUnits{Timeout: config.Duration(5 * time.Second)}
	})
}
