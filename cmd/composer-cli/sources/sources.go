// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package sources

import (
	"github.com/spf13/cobra"

	"github.com/weldr/weldr-client/cmd/composer-cli/root"
)

var (
	sourcesCmd = &cobra.Command{
		Use:   "sources ...",
		Short: "Manage sources",
		Long:  "Manage project sources on the server",
	}
)

func init() {
	root.AddRootCommand(sourcesCmd)
}
