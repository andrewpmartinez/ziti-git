package main

import (
	"github.com/andrewpmartinez/ziti-git/zg"
	"github.com/spf13/cobra"
)

func NewTableStatusCmd(ctx *Ctx) *cobra.Command{
	tableStatusCmd := &cobra.Command{
		Use:     "table-status [-t <tag>]",
		Aliases: []string{"ts"},
		Short:   "show the table status of all the repos or of a specific tag",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			checkRepos(ctx.Repos)
			tag := rootCmd.Flag(FlagTag).Value.String()
			zg.TableStatus(ctx.Repos, tag)
		},
	}
	tableStatusCmd.Flags().StringP(FlagTag, "t", "", "limits actions to repos with <tag>")

	return tableStatusCmd
}
