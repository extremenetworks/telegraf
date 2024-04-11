//go:build !linux

package dmcache

import (
	"github.com/extremenetworks/telegraf"
)

func (*DMCache) Gather(_ telegraf.Accumulator) error {
	return nil
}

func dmSetupStatus() ([]string, error) {
	return []string{}, nil
}
