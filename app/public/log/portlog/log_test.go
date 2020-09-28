package portlog

import "testing"

func TestLog(t *testing.T) {
	LOG.Debug("debug")
	LOG.Info("info")
	LOG.Warning("warning")
}
