package main

import (
	"github.com/andrewpmartinez/ziti-git/zg"
	"github.com/spf13/cobra"
	"path/filepath"
)

func NewUnregisterCmd(ctx *Ctx) *cobra.Command {
	unregisterCmd := &cobra.Command{
		Use:     "unregister <repo>",
		Aliases: []string{"u"},
		Short:   "unregister <repo>",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			checkRepos(ctx.Repos)

			path, _ := filepath.Abs(args[0])
			zg.UnregisterRepo(path, ctx.Repos)
		},
	}

	return unregisterCmd
}
