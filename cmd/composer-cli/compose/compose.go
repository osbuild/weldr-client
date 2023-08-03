// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

package compose

import (
	"github.com/spf13/cobra"

	"github.com/osbuild/weldr-client/v2/cmd/composer-cli/root"
)

var (
	composeCmd = &cobra.Command{
		Use:   "compose ...",
		Short: "Manage composes",
		Long: `Manage composes on the server

  The 'start' and 'start-ostree' commands can optionally upload the results,
  view the full upload profile reference here:
  https://www.osbuild.org/guides/image-builder-on-premises/uploading-to-cloud.html
`,
		Example: `  TOML profile for uploading to AWS

    provider = "aws"

    [settings]
    accessKeyID = "AWS_ACCESS_KEY_ID"
    secretAccessKey = "AWS_SECRET_ACCESS_KEY"
    bucket = "AWS_BUCKET"
    region = "AWS_REGION"
    key = "OBJECT_KEY"

  TOML profile for uploading to GCP

    provider = "gcp"

    [settings]
    bucket = "GCP_BUCKET"
    region = "GCP_STORAGE_REGION"
    object = "OBJECT_KEY"
    credentials = "GCP_CREDENTIALS"

  TOML profile for uploading to Azure

    provider = "azure"

    [settings]
    storageAccount = "your storage account name"
    storageAccessKey = "storage access key you copied in the Azure portal"
    container = "your storage container name"
`,
	}
)

func init() {
	root.AddRootCommand(composeCmd)
}
