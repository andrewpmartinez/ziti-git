package main

import (
	"fmt"
	"github.com/andrewpmartinez/ziti-git/zg"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var rootCmd *cobra.Command


func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	zg.SetConfigFilePath()
	repos := zg.GetRepos()

	rootCmd = &cobra.Command{
		Use:   "ziti-git",
		Short: "Ziti Git is a multi-repo git tool with additions for the open ziti project!",
	}

	registerCmd := &cobra.Command{
		Use:     "register <tag> <path>",
		Aliases: []string{"r"},
		Short:   "add the repo in <path> to the list of repos, with an optional <tag>",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tag := ""
			path := ""
			if len(args) > 1 {
				tag = args[0]
				path = args[1]
			} else {
				path = args[0]
			}

			path, _ = filepath.Abs(path)
			zg.RegisterRepo(path, tag, repos)
		},
	}

	tableStatusCmd := &cobra.Command{
		Use:     "table-status [<tag>]",
		Aliases: []string{"ts"},
		Short:   "show the table status of all the repos or of a specific tag",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			checkRepos(repos)

			tag := ""

			if len(args) > 0 {
				tag = args[0]
			}

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
		Use:     "list [<tag>]",
		Aliases: []string{"l"},
		Short:   "list all repos or repos for <tag>",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			checkRepos(repos)

			tag := ""

			if len(args) > 0 {
				tag = args[0]
			}

			zg.PrintRepos(tag, repos)
		},
	}

	branchCmd := &cobra.Command{
		Use:     "branch [<tag>]",
		Aliases: []string{"b"},
		Short:   "list all repo branches or repos in <tag>",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			checkRepos(repos)

			tag := ""

			if len(args) > 0 {
				tag = args[0]
			}

			zg.PrintRepos(tag, repos)
		},
	}

	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(tableStatusCmd)
	rootCmd.AddCommand(unregisterCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(branchCmd)

}

func checkRepos(repos []zg.Repo) {
	if len(repos) == 0 {
		fmt.Println("No repositories registered. Nothing to do.")
		fmt.Println("Please register a repository with the command:")
		fmt.Println("ziti-git register [path]")
		os.Exit(0)
	}
}
