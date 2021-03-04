package main

import (
	"github.com/andrewpmartinez/ziti-git/zg"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func NewGitCommand(ctx *Ctx) *cobra.Command {
	executeCmd := &cobra.Command{
		Use:                   "git [-t <tag>] <git commands/args>",
		Aliases:               []string{"g"},
		Short:                 "execute git commands across all repositories or specific <tag> repositories",
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagParsing:    true,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			tag := rootCmd.Flag(FlagTag).Value.String()

			var passArgs []string

			//wish we could do this w/ cobra.Command parsing but
			//it doesn't like pass through type commands
			pastThisCommand := false
			readFirstArg := false
			nextIsTag := false
			for i, arg := range os.Args {
				if !pastThisCommand {
					if arg == cmd.Name() || arg == cmd.Aliases[0] {
						pastThisCommand = true
					}
					continue
				}

				if !readFirstArg {
					readFirstArg = true

					if arg == "-t" || arg == "--tag" {
						if len(os.Args) > (i + 1) {
							tag = os.Args[i+1]
							nextIsTag = true
							continue
						} else {
							cmd.PrintErrf("Error: flag needs an argument: '%s' in %s\n", strings.ReplaceAll(arg, "-", ""), arg)
							_ = cmd.Help()
							return
						}
					} else if arg == "-h" || arg == "--help" {
						_ = cmd.Help()
						return
					}
				}

				if nextIsTag {
					nextIsTag = false
					continue
				}

				passArgs = append(passArgs, arg)

			}

			if len(passArgs) == 0 {
				cmd.PrintErr("Error: no git commands/arguments provided\n")
				_ = cmd.Help()
				return
			}

			zg.RunGitCommand(ctx.Repos, tag, passArgs...)
		},
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
	}

	executeCmd.Flags().StringP(FlagTag, "t", "", "limits actions to repos with <tag>")

	return executeCmd
}
