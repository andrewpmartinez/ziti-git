package main

import (
	"github.com/andrewpmartinez/ziti-git/zg"
	"github.com/spf13/cobra"
)

func NewListCmd(ctx *Ctx) *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "list [-t <tag>]",
		Aliases: []string{"l"},
		Short:   "list all repos or repos for <tag>",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			checkRepos(ctx.Repos)

			tag := rootCmd.Flag(FlagTag).Value.String()

			zg.PrintRepos(tag, ctx.Repos)
		},
	}
	listCmd.Flags().StringP(FlagTag, "t", "", "limits actions to repos with <tag>")

	return listCmd
}
