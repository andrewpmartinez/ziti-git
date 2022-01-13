package main

import (
	"github.com/spf13/cobra"
	"os"
)

func NewUseRemoteCmd(_ *Ctx) *cobra.Command {
	useRemoteCmd := &cobra.Command{
		Use:     "use-remote [-h] [-r <repos>]",
		Aliases: []string{"ur"},
		Short:   "short cut for use-local -u",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			workOnCurrent, _ := cmd.Flags().GetBool("current")
			reposToReplace, _ := cmd.Flags().GetStringArray("repos")
			noTidy, _ := cmd.Flags().GetBool("no-tidy")

			workingDir, _ := os.Getwd()
			var repoDirs []string

			alterGoMods(workOnCurrent, true, reposToReplace, workingDir, repoDirs, noTidy)
		},
	}

	useRemoteCmd.Flags().BoolP("current", "c", false, "only alter the current repository, must be in a git repository folder")
	useRemoteCmd.Flags().StringArrayP("repos", "r", []string{`github\.com/openziti/.*`}, "alter specific replace directives by repository URL regexp, may be specified multiple times")
	useRemoteCmd.Flags().BoolP("no-tidy", "n", false, "if specified, go.mod altering commands will not run go mod tidy")

	return useRemoteCmd
}
