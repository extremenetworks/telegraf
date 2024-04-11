//go:build !linux

package cgroup

import (
	"github.com/extremenetworks/telegraf"
)

func (*CGroup) Gather(_ telegraf.Accumulator) error {
	return nil
}
