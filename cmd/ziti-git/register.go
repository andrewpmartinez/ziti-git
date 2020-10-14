package main

import (
	"github.com/andrewpmartinez/ziti-git/zg"
	"github.com/spf13/cobra"
	"path/filepath"
)

func NewRegisterCmd(ctx *Ctx) *cobra.Command {
	registerCmd := &cobra.Command{
		Use:     "register [-t <tag>] <path>",
		Aliases: []string{"r"},
		Short:   "add the repo in <path> to the list of repos, with an optional <tag>",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tag := rootCmd.Flag(FlagTag).Value.String()
			path := args[0]
			path, _ = filepath.Abs(path)
			zg.RegisterRepo(path, tag, ctx.Repos)
		},
	}
	registerCmd.Flags().StringP(FlagTag, "t", "", "limits actions to repos with <tag>")

	return registerCmd
}
