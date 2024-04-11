//go:build windows

package socketstat

import (
	_ "embed"

	"github.com/extremenetworks/telegraf"
	"github.com/extremenetworks/telegraf/plugins/inputs"
)

//go:embed sample.conf
var sampleConfig string

type Socketstat struct {
	Log telegraf.Logger `toml:"-"`
}

func (s *Socketstat) Init() error {
	s.Log.Warn("current platform is not supported")
	return nil
}
func (*Socketstat) SampleConfig() string                { return sampleConfig }
func (*Socketstat) Gather(_ telegraf.Accumulator) error { return nil }

func init() {
	inputs.Add("socketstat", func() telegraf.Input {
		return &Socketstat{}
	})
}
