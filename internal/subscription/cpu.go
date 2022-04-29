package subscription

import (
	"fmt"
	"log"
	"os"
)

type SubscriptionConfigrationStatsCpu struct {
	ConfigDir string
}

var cfgFileName string = "cpu.conf"
var cfgStatsCpuTemplate string = `
# Read metrics about cpu usage
[[inputs.cpu]]
  ## Whether to report per-cpu stats or not
  percpu = true
  ## Whether to report total system cpu stats or not
  totalcpu = true
  ## If true, collect raw CPU time metrics
  collect_cpu_time = false
  ## If true, compute and report the sum of all non-idle CPU states
  report_active = false
  interval = "%ds"
`

func (h *SubscriptionConfigrationStatsCpu) UpdateConfig(cfg SubscriptionRequest) error {
	cfgFullName := h.ConfigDir + "/" + cfgFileName
	cfgFileContent := fmt.Sprintf(cfgStatsCpuTemplate, cfg.Interval)
	err := os.WriteFile(cfgFullName, []byte(cfgFileContent), 0644)
	if err == nil {
		log.Printf("I! Updated config file %s\n", cfgFullName)
	} else {
		log.Printf("E! Failed to write file %s: %+v\n", cfgFullName, err)
	}
	return err
}

func (h *SubscriptionConfigrationStatsCpu) DeleteConfig() error {
	cfgFullName := h.ConfigDir + "/" + cfgFileName
	err := os.Remove(cfgFullName)
	if err == nil {
		log.Printf("I! Removed config file %s\n", cfgFullName)
	} else {
		log.Printf("E! Failed to remove file %s: %+v\n", cfgFullName, err)
	}
	return err
}
