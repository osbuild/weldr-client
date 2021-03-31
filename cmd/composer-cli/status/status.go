// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package status

import (
	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/cmd/composer-cli/root"
)

var (
	statusCmd = &cobra.Command{
		Use:   "status ...",
		Short: "API server status",
		Long:  "API server status",
	}
)

func init() {
	root.AddRootCommand(statusCmd)
}
