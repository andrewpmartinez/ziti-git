package main

import (
	"fmt"
	"github.com/andrewpmartinez/ziti-git/zg"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var rootCmd *cobra.Command

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}

const (
	FlagTag = "tag"
)

func init() {
	zg.SetConfigFilePath()
	repos := zg.GetRepos()

	rootCmd = &cobra.Command{
		Use:   "ziti-git",
		Short: "Ziti Git is a multi-repo git tool with additions for the open ziti project!",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tag := rootCmd.Flag(FlagTag).Value.String()

			passArgs := os.Args[1:]

			if tag != "" {
				passArgs = os.Args[3:]
			}

			zg.RunCmd(repos, tag, passArgs...)
		},
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
	}

	rootCmd.PersistentFlags().StringP(FlagTag, "t", "", "limits actions to repos with <tag>")

	registerCmd := &cobra.Command{
		Use:     "register [-t <tag>] <path>",
		Aliases: []string{"r"},
		Short:   "add the repo in <path> to the list of repos, with an optional <tag>",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tag := rootCmd.Flag(FlagTag).Value.String()
			path := args[0]
			path, _ = filepath.Abs(path)
			zg.RegisterRepo(path, tag, repos)
		},
	}

	tableStatusCmd := &cobra.Command{
		Use:     "table-status [-t <tag>]",
		Aliases: []string{"ts"},
		Short:   "show the table status of all the repos or of a specific tag",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			checkRepos(repos)
			tag := rootCmd.Flag(FlagTag).Value.String()
			zg.TableStatus(repos, tag)
		},
	}

	unregisterCmd := &cobra.Command{
		Use:     "unregister <repo>",
		Aliases: []string{"u"},
		Short:   "unregister <repo>",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			checkRepos(repos)

			path, _ := filepath.Abs(args[0])
			zg.UnregisterRepo(path, repos)
		},
	}

	listCmd := &cobra.Command{
		Use:     "list [-t <tag>]",
		Aliases: []string{"l"},
		Short:   "list all repos or repos for <tag>",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			checkRepos(repos)

			tag := rootCmd.Flag(FlagTag).Value.String()

			zg.PrintRepos(tag, repos)
		},
	}

	branchCmd := &cobra.Command{
		Use:     "branch [-t <tag>]",
		Aliases: []string{"b"},
		Short:   "list all repo branches or repos in <tag>",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			checkRepos(repos)
			tag := rootCmd.Flag(FlagTag).Value.String()

			zg.PrintRepos(tag, repos)
		},
	}

	clone := &cobra.Command{
		Use:   "clone",
		Short: "clones the core openziti repos to the current directory",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {

			for _, repo := range []string{
				"git@github.com:openziti/edge.git",
				"git@github.com:openziti/fabric.git",
				"git@github.com:openziti/foundation.git",
				"git@github.com:openziti/ziti.git",
				"git@github.com:openziti/sdk-golang.git",
			} {
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
					continue
				}

				println(stdout.Bytes())
			}

		},
	}

	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(tableStatusCmd)
	rootCmd.AddCommand(unregisterCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(branchCmd)
	rootCmd.AddCommand(clone)
}

func checkRepos(repos []zg.Repo) {
	if len(repos) == 0 {
		fmt.Println("No repositories registered. Nothing to do.")
		fmt.Println("Please register a repository with the command:")
		fmt.Println("ziti-git register [path]")
		os.Exit(0)
	}
}
