// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package projects

import (
	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

var (
	projectsCmd = &cobra.Command{
		Use:   "projects ...",
		Short: "Project related commands",
		Long:  "Project related commands",
	}
)

func init() {
	root.AddRootCommand(projectsCmd)
}
