package opcua

import (
	"github.com/extremenetworks/telegraf"
)

// DebugLogger logs messages from opcua at the debug level.
type DebugLogger struct {
	Log telegraf.Logger
}

func (l *DebugLogger) Write(p []byte) (n int, err error) {
	l.Log.Debug(string(p))
	return len(p), nil
}
