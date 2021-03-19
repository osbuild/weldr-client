// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package modules

import (
	"github.com/spf13/cobra"

	"github.com/weldr/weldr-client/cmd/composer-cli/root"
)

var (
	modulesCmd = &cobra.Command{
		Use:   "modules ...",
		Short: "Module related commands",
		Long:  "Module related commands",
	}
)

func init() {
	root.AddRootCommand(modulesCmd)
}
