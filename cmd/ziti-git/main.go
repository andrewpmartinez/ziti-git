package main

import (
	"fmt"
	"github.com/andrewpmartinez/ziti-git/zg"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var rootCmd *cobra.Command

func main() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}

const (
	FlagTag = "tag"
)

var ZitiRepos = map[string]string{
	"edge":       "git@github.com:openziti/edge.git",
	"fabric":     "git@github.com:openziti/fabric.git",
	"foundation": "git@github.com:openziti/foundation.git",
	"ziti":       "git@github.com:openziti/ziti.git",
	"sdk-golang": "git@github.com:openziti/sdk-golang.git",
}

var ZitiModules = map[string]string{
	"github.com/openziti/edge":       "edge",
	"github.com/openziti/fabric":     "fabric",
	"github.com/openziti/foundation": "foundation",
	"github.com/openziti/sdk-golang": "sdk-golang",
}

func init() {
	zg.SetConfigFilePath()
	repos := zg.GetRepos()

	rootCmd = &cobra.Command{
		Use:   "ziti-git",
		Short: "Ziti Git is a multi-repo git tool with additions for the open ziti project!",
		Args:  cobra.MinimumNArgs(1),
	}
	rootCmd.Flags().StringP(FlagTag, "t", "", "limits actions to repos with <tag>")

	executeCmd := &cobra.Command{
		Use:                   "exec [-t <tag>] <git commands/args>",
		Aliases:               []string{"e"},
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

			zg.RunCmd(repos, tag, passArgs...)
		},
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
	}
	executeCmd.Flags().StringP(FlagTag, "t", "", "limits actions to repos with <tag>")

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
	registerCmd.Flags().StringP(FlagTag, "t", "", "limits actions to repos with <tag>")

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
	tableStatusCmd.Flags().StringP(FlagTag, "t", "", "limits actions to repos with <tag>")

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

	unregisterTagCmd := &cobra.Command{
		Use:     "unregister-tag <tag>",
		Aliases: []string{"ut"},
		Short:   "unregister-tag <tag>",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			checkRepos(repos)
			tag := args[0]

			for _, repo := range repos {
				if repo.Tag == tag {
					println("...unregister: " + repo.Location)
					zg.UnregisterRepo(repo.Location, zg.GetRepos())
				}
			}
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
	listCmd.Flags().StringP(FlagTag, "t", "", "limits actions to repos with <tag>")

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
	branchCmd.Flags().StringP(FlagTag, "t", "", "limits actions to repos with <tag>")

	useLocalCmd := &cobra.Command{
		Use:     "use-local [-hu] [-r <repos>]",
		Aliases: []string{"ul"},
		Short:   "alter go.mod files for ziti repos to use local repositories via replace directives",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			workOnCurrent, _ := cmd.Flags().GetBool("current")
			undo, _ := cmd.Flags().GetBool("undo")
			reposToReplace, _ := cmd.Flags().GetStringArray("repos")

			workingDir, _ := os.Getwd()
			var repoDirs []string

			//fill repoDirs with directories to work on
			if workOnCurrent {
				if !zg.IsPathGitRepo(workingDir) {
					formattedErrorExit("the current directory is not a Git repository [%s]", workingDir)
				}

				repoDirs = append(repoDirs, workingDir)
			} else {

				if zg.IsPathGitRepo(workingDir) {
					var err error
					workingDir, err = filepath.Abs(workingDir + string(os.PathSeparator) + "..")

					if err != nil {
						formattedErrorExit("detected git directory, tried moving to parent but failed: %s", err)
					}

					fmt.Printf("detected git directory, setting working directory to parent: %s\n", workingDir)

				}
				allDirs, err := ioutil.ReadDir(workingDir)

				if err != nil {
					formattedErrorExit("Could not read working directory: %v", err)
				}

				for _, dir := range allDirs {
					if dir.IsDir() && !strings.HasPrefix(".", dir.Name()) {
						path := filepath.Join(workingDir, dir.Name())
						if zg.IsPathGitRepo(path) {
							repoDirs = append(repoDirs, path)
						}
					}
				}

				if len(repoDirs) == 0 {
					formattedErrorExit("the current directory doesn't appear to have any git repositories [%s]", workingDir)
				}
			}

			for _, repoDir := range repoDirs {

				if zg.HasGoModFile(repoDir) {
					if undo {
						fmt.Printf("disabling replacements on: %s\n", repoDir)
						if err := zg.DisableGoModReplaceDirectives(repoDir, reposToReplace); err != nil {
							formattedErrorExit("Could not disable go mod replacements in [%s]: %v", repoDir, err)
						}
					} else {
						fmt.Printf("enabling replacements on: %s\n", repoDir)
						if err := zg.EnableGoModReplaceDirectives(repoDir, reposToReplace); err != nil {
							formattedErrorExit("Could not enable go mod replacements in [%s]: %v", repoDir, err)
						}
					}
				}
			}
		},
	}

	useLocalCmd.Flags().BoolP("current", "c", false, "only alter the current repository, must be in a git repository folder")
	useLocalCmd.Flags().BoolP("undo", "u", false, "alter go.mod files to not use local repositories, may be combined with -h")
	useLocalCmd.Flags().StringArrayP("repos", "r", []string{`github\.com/openziti/.*`}, "alter specific replace directives by repository URL regexp, may be specified multiple times")

	checkoutCmd := &cobra.Command{
		Use:     "checkout",
		Aliases: []string{"co"},
		Short:   "inspects the go.mod file of the openziti/ziti repo to produce a script to checkout exact openziti dependencies necessary",
		Run: func(cmd *cobra.Command, args []string) {
			repoDir, _ := os.Getwd()

			if !zg.IsPathGitRepo(repoDir) {
				formattedErrorExit("current directory is not a git repo [%s]", repoDir)
			}

			info, err := zg.GetGoModInfo(repoDir)

			if err != nil {
				formattedErrorExit("could not get go mod info for repo [%s]: %v", repoDir, err)
			}

			parentPath := repoDir + string(filepath.Separator) + ".."
			absParentPath, err := filepath.Abs(parentPath)

			if err != nil {
				formattedErrorExit("cannot determine parent path for [%s]: %v", parentPath, err)
			}

			script := "\n"
			script += fmt.Sprintf(`cd "%s"`, absParentPath) + "\n"
			for _, require := range info.Require {
				if dirName, found := ZitiModules[require.Mod.Path]; found {
					script += fmt.Sprintf(`git -C "%s" checkout %s`, "./"+dirName, require.Mod.Version) + "\n"
				}
			}

			println(script)
		},
	}

	clone := &cobra.Command{
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
	clone.Flags().StringP(FlagTag, "t", "", "adds cloned repositories to <tag>")
	clone.Flags().BoolP("register", "r", false, "add cloned repos to ziti-git under <tag> if specified")

	rootCmd.AddCommand(executeCmd)
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(tableStatusCmd)
	rootCmd.AddCommand(unregisterCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(branchCmd)
	rootCmd.AddCommand(clone)
	rootCmd.AddCommand(unregisterTagCmd)
	rootCmd.AddCommand(useLocalCmd)
	rootCmd.AddCommand(checkoutCmd)
}

func checkRepos(repos []zg.Repo) {
	if len(repos) == 0 {
		fmt.Println("No repositories registered. Nothing to do.")
		fmt.Println("Please register a repository with the command:")
		fmt.Println("ziti-git register [path]")
		os.Exit(0)
	}
}

func errorExitWithCode(msg string, exitCode int) {
	_, _ = os.Stderr.WriteString("error: " + msg + "\n")
	os.Exit(exitCode)
}

func formattedErrorExit(msg string, formatArgs ...interface{}) {
	msg = fmt.Sprintf(msg, formatArgs...)
	errorExitWithCode(msg, 1)
}
