package main

import (
	"github.com/andrewpmartinez/ziti-git/zg"
	"github.com/spf13/cobra"
)

func NewUnregisterTagCmd(ctx *Ctx) *cobra.Command {
	return &cobra.Command{
		Use:     "unregister-tag <tag>",
		Aliases: []string{"ut"},
		Short:   "unregister-tag <tag>",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			checkRepos(ctx.Repos)
			tag := args[0]

			for _, repo := range ctx.Repos {
				if repo.Tag == tag {
					println("...unregister: " + repo.Location)
					zg.UnregisterRepo(repo.Location, zg.GetRepos())
				}
			}
		},
	}
}
