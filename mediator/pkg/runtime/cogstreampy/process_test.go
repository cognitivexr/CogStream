package cogstreampy

import "testing"

func Test_parseEngineAddressFromLog(t *testing.T) {
	line := "INFO:cogstream.engine.srv:started server socket on address ('0.0.0.0', 46699)"

	addr, _ := parseEngineAddressFromLog(line)
	if addr != "0.0.0.0:46699" {
		t.Errorf("invalid address %s", addr)
	}

}
