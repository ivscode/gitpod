// Copyright (c) 2022 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package cmd

import (
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/gitpod-io/gitpod/previewctl/pkg/preview"
	"github.com/spf13/cobra"
)

var (
	watch = false
)

func installContextCmd(logger log.Logger) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "install-context",
		Short: "Installs the kubectl context of a preview environment.",
		Run: func(cmd *cobra.Command, args []string) {
			p := preview.New(branch, logger)

			err := p.InstallContext(watch)
			if err != nil {
				level.Error(logger).Log("msg", "Failed to install context", "branch", p.Branch, "err", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().BoolVar(&watch, "watch", false, "If wait is enabled, previewctl will keep trying to install the kube-context every 30 seconds.")
	return cmd
}
