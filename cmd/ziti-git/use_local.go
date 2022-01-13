package main

import (
	"fmt"
	"github.com/andrewpmartinez/ziti-git/zg"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func NewUseLocalCmd(_ *Ctx) *cobra.Command {
	useLocalCmd := &cobra.Command{
		Use:     "use-local [-hu] [-r <repos>]",
		Aliases: []string{"ul"},
		Short:   "alter go.mod files for ziti repos to use local repositories via replace directives",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			workOnCurrent, _ := cmd.Flags().GetBool("current")
			undo, _ := cmd.Flags().GetBool("undo")
			reposToReplace, _ := cmd.Flags().GetStringArray("repos")
			noTidy, _ := cmd.Flags().GetBool("no-tidy")

			workingDir, _ := os.Getwd()
			var repoDirs []string

			alterGoMods(workOnCurrent, undo, reposToReplace, workingDir, repoDirs, noTidy)
		},
	}

	useLocalCmd.Flags().BoolP("current", "c", false, "only alter the current repository, must be in a git repository folder")
	useLocalCmd.Flags().BoolP("undo", "u", false, "alter go.mod files to not use local repositories, may be combined with -h")
	useLocalCmd.Flags().StringArrayP("repos", "r", []string{`github\.com/openziti/.*`}, "alter specific replace directives by repository URL regexp, may be specified multiple times")
	useLocalCmd.Flags().BoolP("no-tidy", "n", false, "if specified, go.mod altering commands will not run go mod tidy")

	return useLocalCmd
}

func alterGoMods(workOnCurrent bool, undo bool, reposToReplace []string, workingDir string, repoDirs []string, noTidy bool) {
	//fill repoDirs with directories to work on
	if workOnCurrent {
		if !zg.IsPathGitRepo(workingDir) {
			formattedErrorExit("the current directory is not a Git repository [%s]", workingDir)
		}

		repoDirs = append(repoDirs, workingDir)
	} else {

		if zg.IsPathGitRepo(workingDir) {
			var err error
			workingDir, err = filepath.Abs(workingDir + string(os.PathSeparator) + "..")

			if err != nil {
				formattedErrorExit("detected git directory, tried moving to parent but failed: %s", err)
			}

			fmt.Printf("detected git directory, setting working directory to parent: %s\n", workingDir)

		}
		allDirs, err := ioutil.ReadDir(workingDir)

		if err != nil {
			formattedErrorExit("Could not read working directory: %v", err)
		}

		for _, dir := range allDirs {
			if dir.IsDir() && !strings.HasPrefix(".", dir.Name()) {
				path := filepath.Join(workingDir, dir.Name())
				if zg.IsPathGitRepo(path) {
					repoDirs = append(repoDirs, path)
				}
			}
		}

		if len(repoDirs) == 0 {
			formattedErrorExit("the current directory doesn't appear to have any git repositories [%s]", workingDir)
		}
	}

	for _, repoDir := range repoDirs {

		if zg.HasGoModFile(repoDir) {
			if undo {
				fmt.Printf("disabling replacements on: %s\n", repoDir)
				if err := zg.DisableGoModReplaceDirectives(repoDir, reposToReplace); err != nil {
					formattedErrorExit("Could not disable go mod replacements in [%s]: %v", repoDir, err)
				}

				if !noTidy {
					fmt.Printf("- tidying: %s\n", repoDir)
					zg.GoModTidy(repoDir)
				}

			} else {
				fmt.Printf("enabling replacements on: %s\n", repoDir)
				if err := zg.EnableGoModReplaceDirectives(repoDir, reposToReplace); err != nil {
					formattedErrorExit("Could not enable go mod replacements in [%s]: %v", repoDir, err)
				}

				if !noTidy {
					fmt.Printf("- tidying: %s\n", repoDir)
					zg.GoModTidy(repoDir)
				}
			}
		}
	}
}
