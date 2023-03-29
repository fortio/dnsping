package main

import (
	"os"
	"testing"

	"fortio.org/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"dnsping": Main,
	}))
}

func TestDNSPing(t *testing.T) {
	testscript.Run(t, testscript.Params{Dir: "./"})
}
