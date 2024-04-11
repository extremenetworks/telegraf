//go:build !linux

package infiniband

import (
	"github.com/extremenetworks/telegraf"
	"github.com/extremenetworks/telegraf/plugins/inputs"
)

func (i *Infiniband) Init() error {
	i.Log.Warn("Current platform is not supported")
	return nil
}

func (*Infiniband) Gather(_ telegraf.Accumulator) error {
	return nil
}

func init() {
	inputs.Add("infiniband", func() telegraf.Input {
		return &Infiniband{}
	})
}
