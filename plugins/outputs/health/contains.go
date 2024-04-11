package health

import "github.com/extremenetworks/telegraf"

type Contains struct {
	Field string `toml:"field"`
}

func (c *Contains) Check(metrics []telegraf.Metric) bool {
	success := false
	for _, m := range metrics {
		ok := m.HasField(c.Field)
		if ok {
			success = true
		}
	}

	return success
}
