//go:build !custom || processors || processors.aws_ec2

package all

import _ "github.com/extremenetworks/telegraf/plugins/processors/aws_ec2" // register plugin
