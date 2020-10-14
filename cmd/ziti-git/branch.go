package main

import (
	"github.com/andrewpmartinez/ziti-git/zg"
	"github.com/spf13/cobra"
)

func NewBranchCmd(ctx *Ctx) *cobra.Command {
	branchCmd := &cobra.Command{
		Use:     "branch [-t <tag>]",
		Aliases: []string{"b"},
		Short:   "list all repo branches or repos in <tag>",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			checkRepos(ctx.Repos)
			tag := rootCmd.Flag(FlagTag).Value.String()

			zg.PrintRepos(tag, ctx.Repos)
		},
	}
	branchCmd.Flags().StringP(FlagTag, "t", "", "limits actions to repos with <tag>")

	return branchCmd
}
