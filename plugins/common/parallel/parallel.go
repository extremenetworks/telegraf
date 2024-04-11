package parallel

import "github.com/extremenetworks/telegraf"

type Parallel interface {
	Enqueue(telegraf.Metric)
	Stop()
}
