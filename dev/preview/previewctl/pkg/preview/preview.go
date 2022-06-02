// Copyright (c) 2022 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package preview

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type Preview struct {
	Branch string

	logger log.Logger
}

func New(branch string, logger log.Logger) *Preview {
	if branch == "" {
		out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
		if err != nil {
			level.Error(logger).Log("msg", "Could not retrieve branch name", "err", err)
			os.Exit(1)
		}
		branch = string(out)
	} else {
		_, err := exec.Command("git", "rev-parse", "--verify", branch).Output()
		if err != nil {
			level.Error(logger).Log("msg", "Branch does not exist", "branch", branch)
			os.Exit(1)
		}
	}

	return &Preview{
		Branch: branch,
		logger: logger,
	}
}

func (p *Preview) InstallContext(watch bool) error {
	if watch {
		// The most precise implementation for a watch-loop would be to implement
		// a pub-sub like logic where previewctl would react to changes of a preview IP address.
		// For now just an infinite loop will do!
		installTicker := time.NewTicker(30 * time.Second)

		// nolint:gosimple
		for {
			select {
			case <-installTicker.C:
				if err := installContext(p.Branch); err == nil {
					// No error means successful context installation
					return nil
				}
			}
		}
	}

	return installContext(p.Branch)
}

func installContext(branch string) error {
	return exec.Command("bash", "/workspace/gitpod/dev/preview/install-k3s-kubeconfig.sh", "-b", branch).Run()
}

func (p *Preview) GetPreviewName() string {
	withoutRefsHead := strings.Replace(p.Branch, "/refs/heads/", "", 1)
	lowerCased := strings.ToLower(withoutRefsHead)

	var re = regexp.MustCompile(`[^-a-z0-9]`)
	sanitizedBranch := re.ReplaceAllString(lowerCased, `$1-$2`)

	if len(sanitizedBranch) > 20 {
		h := sha256.New()
		h.Write([]byte(sanitizedBranch))
		hashedBranch := hex.EncodeToString(h.Sum(nil))

		sanitizedBranch = sanitizedBranch[0:10] + hashedBranch[0:10]
	}

	return sanitizedBranch
}
