//go:build windows

package processes

import (
	"github.com/extremenetworks/telegraf"
	"github.com/extremenetworks/telegraf/plugins/inputs"
)

type Processes struct {
	Log telegraf.Logger
}

func (e *Processes) Init() error {
	e.Log.Warn("Current platform is not supported")
	return nil
}

func (e *Processes) Gather(_ telegraf.Accumulator) error {
	return nil
}

func init() {
	inputs.Add("processes", func() telegraf.Input {
		return &Processes{}
	})
}
