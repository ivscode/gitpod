// Copyright (c) 2022 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package main

import (
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/gitpod-io/gitpod/previewctl/cmd"
)

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	root := cmd.RootCmd(logger)
	if err := root.Execute(); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}
}
