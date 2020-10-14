package main

import (
	"fmt"
	"github.com/andrewpmartinez/ziti-git/zg"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func NewCheckoutCmd(_ *Ctx) *cobra.Command {
	checkoutCmd := &cobra.Command{
		Use:     "checkout",
		Aliases: []string{"co"},
		Short:   "inspects the go.mod file of the openziti/ziti repo to produce a script to checkout exact openziti dependencies necessary",
		Run: func(cmd *cobra.Command, args []string) {
			repoDir, _ := os.Getwd()

			if !zg.IsPathGitRepo(repoDir) {
				formattedErrorExit("current directory is not a git repo [%s]", repoDir)
			}

			info, err := zg.GetGoModInfo(repoDir)

			if err != nil {
				formattedErrorExit("could not get go mod info for repo [%s]: %v", repoDir, err)
			}

			parentPath := repoDir + string(filepath.Separator) + ".."
			absParentPath, err := filepath.Abs(parentPath)

			if err != nil {
				formattedErrorExit("cannot determine parent path for [%s]: %v", parentPath, err)
			}

			script := "\n"
			script += fmt.Sprintf(`cd "%s"`, absParentPath) + "\n"
			for _, require := range info.Require {
				if dirName, found := ZitiModules[require.Mod.Path]; found {
					script += fmt.Sprintf(`git -C "%s" checkout %s`, "./"+dirName, require.Mod.Version) + "\n"
				}
			}

			println(script)
		},
	}

	return checkoutCmd
}
