package main

import (
	"fmt"
	"github.com/andrewpmartinez/ziti-git/zg"
	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
	"io"
	"os"
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
	"edge":         "git@github.com:openziti/edge.git",
	"fabric":       "git@github.com:openziti/fabric.git",
	"foundation":   "git@github.com:openziti/foundation.git",
	"ziti":         "git@github.com:openziti/ziti.git",
	"sdk-golang":   "git@github.com:openziti/sdk-golang.git",
	"sdk-js":       "git@github.com:openziti/ziti-sdk-js.git",
	"sdk-jvm":      "git@github.com:openziti/ziti-sdk-jvm.git",
	"doc":          "git@github.com:openziti/openziti.github.io.git",
	"tunnel-sdk-c": "git@github.com:openziti/ziti-tunnel-sdk-c.git",
	"channel":      "git@github.com:openziti/channel.git",
}

var ZitiModules = map[string]string{
	"github.com/openziti/edge":       "edge",
	"github.com/openziti/fabric":     "fabric",
	"github.com/openziti/foundation": "foundation",
	"github.com/openziti/sdk-golang": "sdk-golang",
	"github.com/openziti/channel":    "channel",
}

type Ctx struct {
	Repos  []zg.Repo
	Output io.Writer
}

func init() {
	zg.SetConfigFilePath()
	ctx := &Ctx{Repos: zg.GetRepos(), Output: colorable.NewColorableStdout()}

	rootCmd = &cobra.Command{
		Use:   "ziti-git",
		Short: "Ziti Git is a multi-repo git tool with additions for the open ziti project!",
		Args:  cobra.MinimumNArgs(1),
	}
	rootCmd.Flags().StringP(FlagTag, "t", "", "limits actions to repos with <tag>")

	subCmds := []*cobra.Command{
		NewGitCommand(ctx),
		NewExecuteCommand(ctx),
		NewRegisterCmd(ctx),
		NewTableStatusCmd(ctx),
		NewUnregisterCmd(ctx),
		NewUnregisterTagCmd(ctx),
		NewListCmd(ctx),
		NewBranchCmd(ctx),
		NewUseLocalCmd(ctx),
		NewCheckoutCmd(ctx),
		NewCloneCmd(ctx),
		NewUseRemoteCmd(ctx),
	}

	for _, subCmd := range subCmds {
		rootCmd.AddCommand(subCmd)
	}

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
