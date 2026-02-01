// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Hubble

package main

import (
	"fmt"
	"os"

	"github.com/cilium/cilium/hubble/cmd"
	"github.com/google/gops/agent"
)

func main() {
	// Start the gops agent with default options.
	if err := agent.Listen(agent.Options{}); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
