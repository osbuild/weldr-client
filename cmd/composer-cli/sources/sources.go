// Copyright 2021 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package sources

import (
	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

var (
	sourcesCmd = &cobra.Command{
		Use:   "sources ...",
		Short: "Manage sources",
		Long: `Manage project sources on the server

  The full source reference can be found here:
  https://www.osbuild.org/guides/user-guide/managing-repositories.html
`,
		Example: `  TOML source for 3rd party rpm repository without gpg checking

    id = "extra-repo"
    name = "Extra rpm repository"
    type = "yum-baseurl"
    url = "https://repo.nowhere.com/extra/"
    check_gpg = false
`,
	}
)

func init() {
	root.AddRootCommand(sourcesCmd)
}
