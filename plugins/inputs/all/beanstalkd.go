//go:build !custom || inputs || inputs.beanstalkd

package all

import _ "github.com/extremenetworks/telegraf/plugins/inputs/beanstalkd" // register plugin
