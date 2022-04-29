# Subscription server
Subscription server implements Alert API subscription. See OpenAPI specification
for mor details

https://github.com/extremenetworks/wns_openapi/blob/main/openapi/alerts/platform.yaml

# How to use subscription server
Start `telegraf` with parameters `--subscription-addr` and `--config-directory`. Config directory must exists and can be empty. Specify a port for `--subscription-addr`, that is not used by other applications:
```
telegraf --config telegraf.conf --subscription-addr :35000 --config-directory conf.d
```
Send subscription command using curl:
```
curl -X POST http://localhost:35000/v1/subscribe/stats/cpu -d '{ "interval": 30, "samplePeriod": 6 }'
```
Telegraf will create/update config file:
```
% cat ~/git/telegraf/conf.d/cpu.conf                                                                   

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
  interval = "30s"
```

Telegraf will re-read configuration and restart automatically:
```
2022-04-29T19:49:53Z I! Updated config file conf.d/cpu.conf
2022-04-29T19:49:53Z I! Reloading Telegraf config
2022-04-29T19:49:53Z I! [agent] Hang on, flushing any cached metrics before shutdown
2022-04-29T19:49:53Z I! [agent] Stopping running outputs
2022-04-29T19:49:53Z I! Starting Telegraf 1.21.4-55d994a6
```

# Structure of the code
* `subscription.go` - Starts HTTP listeners, instantiates stats/state handlers, perform generic handling of REST API requests. Stats/State handlers must implement `SubscriptionConfigration` interface.
* `cpu.go` - Updates configuration for `cpu` input plugin. Updates/removes config file `cpu.conf`.

# Future development
Use `cpu.go` as an example when adding handlers for updating plugings for stats and state updates.
For each new handler, the following steps must be implemented:
* Add a source file to package `subscription`, named according to the metric, for example `mem.go` or `temperature.go`
* Define config file name for the metric. The config file name should be unique for the metric. Example: `var cfgFileName string = "cpu.conf"`
* Define template for the content of the config file, that corresponds to the plugin. Example:
```
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
```
* Implement methods
    * `UpdateConfig(cfg SubscriptionRequest) error`
    * `DeleteConfig() error`
* Update file `subscrition.go` method `Start()` by adding new handler to map of handlers. Example: `h.statsConfigHandlers["cpu"] = &SubscriptionConfigrationStatsCpu{ConfigDir: h.ConfigDir}`
In the map the key is metric name. The metric name must match last part of the REST API endpoint path. For example, path for subscription to CPU utilization is `v1/subscribe/stats/cpu`, so name of the metric is `cpu`.

