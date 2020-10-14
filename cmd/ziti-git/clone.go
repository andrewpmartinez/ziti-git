package main

import (
	"fmt"
	"github.com/andrewpmartinez/ziti-git/zg"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"path/filepath"
)

func NewCloneCmd(_ *Ctx) *cobra.Command {
	cloneCmd := &cobra.Command{
		Use:     "clone [-t <tag>] [-r]",
		Aliases: []string{"c"},
		Short:   "clones the core openziti repos to the current directory",
		Run: func(cmd *cobra.Command, _ []string) {

			for dir, repo := range ZitiRepos {
				tag := rootCmd.Flag(FlagTag).Value.String()

				color.Cyan("Cloning: %s", repo)
				stdout, stderr, err := zg.Exec("git", "clone", repo)

				if err != nil {
					_, _ = color.Error.Write([]byte(err.Error()))
					if stderr != nil && stderr.Len() != 0 {
						_, _ = color.Error.Write(stderr.Bytes())
					}
					continue
				}

				if stderr != nil && stderr.Len() != 0 {
					_, _ = color.Error.Write(stderr.Bytes())
				}
				if stdout != nil && stdout.Len() != 0 {
					println(stdout.Bytes())
				}

				if cmd.Flag("register").Value.String() == "true" {
					repoPath, _ := filepath.Abs(filepath.Join(".", dir))
					fmt.Printf("...registering as %s -> %s\n", tag, repoPath)
					zg.RegisterRepo(repoPath, tag, zg.GetRepos())
				}
			}
		},
	}
	cloneCmd.Flags().StringP(FlagTag, "t", "", "adds cloned repositories to <tag>")
	cloneCmd.Flags().BoolP("register", "r", false, "add cloned repos to ziti-git under <tag> if specified")

	return cloneCmd
}
