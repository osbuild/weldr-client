// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"github.com/spf13/cobra"

	"github.com/weldr/weldr-client/cmd/composer-cli/root"
)

var (
	blueprintsCmd = &cobra.Command{
		Use:   "blueprints ...",
		Short: "Manage blueprints",
		Long:  "Manage blueprints on the server",
	}
)

func init() {
	root.AddRootCommand(blueprintsCmd)
}
