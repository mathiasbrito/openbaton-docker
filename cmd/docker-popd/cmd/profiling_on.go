// +build profiling

package cmd

import (
	// plug in http live profiling support
	_ "net/http/pprof"
)

const profilingSupport = true