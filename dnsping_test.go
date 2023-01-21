package main

import (
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"dnsping": Main,
	}))
}

func TestDelta(t *testing.T) {
	testscript.Run(t, testscript.Params{Dir: "./"})
}
