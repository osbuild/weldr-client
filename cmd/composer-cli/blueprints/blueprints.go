// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package blueprints

import (
	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

var (
	blueprintsCmd = &cobra.Command{
		Use:   "blueprints ...",
		Short: "Manage blueprints",
		Long: `Manage blueprints on the server

  The full blueprint reference can be found here:
  https://www.osbuild.org/guides/image-builder-on-premises/blueprint-reference.html
`,
		Example: `  TOML blueprint for an image with tmux, a user, and a new group for the user:

    name = "student-image"
    description = "A base system with a student account"
    version = "0.0.1"

    [[packages]]
    name = "tmux"
    version = "*"

    [[customizations.user]]
    name = "bart"
    description = "Student account for Bart"
    groups = ["students"]
    password = "$6$CHO2$3rN8eviE2t50lmVyBYihTgVRHcaecmeCk31LeOUleVK/R/aeWVHVZDi26zAH.o0ywBKH9Tc0/wm7sW/q39uyd1"

    [[customizations.group]]
    name = "students"


  TOML blueprint for an image with a custom filesystem:

    name = "custom-fs-image"
    description = "A base system with a custom filesystem"
    version = "0.0.1"

    [[packages]]
    name = "tmux"
    version = "*"

    [[customizations.filesystem]]
    mountpoint = "/"
    size = "2GiB"

    [[customizations.filesystem]]
    mountpoint = "/home"
    size = "5GiB"
`,
	}
)

func init() {
	root.AddRootCommand(blueprintsCmd)
}
