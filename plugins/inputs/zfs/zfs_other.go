//go:build !linux && !freebsd

package zfs

import (
	"github.com/extremenetworks/telegraf"
	"github.com/extremenetworks/telegraf/plugins/inputs"
)

func (*Zfs) Gather(_ telegraf.Accumulator) error {
	return nil
}

func init() {
	inputs.Add("zfs", func() telegraf.Input {
		return &Zfs{}
	})
}
