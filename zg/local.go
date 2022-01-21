package zg

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

// knownRepos is a list of shortcuts for use-local and use-remote
// so that full regular expressions don't always have to be used
// when altering go.mod files. RegularExpressions are still the
// input to deal with any repo that may be part of the project
var knownRepos = map[string]string{
	"edge":       "github.com/openziti/edge.*",
	"fabric":     "github.com/openziti/fabric.*",
	"foundation": "github.com/openziti/foundation.*",
	"sdk-golang": "github.com/openziti/sdk-golang.*",
	"sdk":        "github.com/openziti/sdk-golang.*",
}

func IsPathGitRepo(path string) bool {
	dirInfo, _ := os.Stat(path)
	if dirInfo == nil || !dirInfo.IsDir() {
		return false
	}

	dirs, _ := ioutil.ReadDir(path)

	for _, dir := range dirs {
		if dir.Name() == ".git" {
			return true
		}
	}

	return false
}

func GoModTidy(repoDir string) {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = repoDir

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	color.Cyan(repoDir)
	err := cmd.Run()

	hasOutput := false
	if err != nil {
		hasOutput = true
		_, _ = fmt.Fprint(os.Stderr, stderr.String())
	}

	outStdString := out.String()

	if !hasOutput && outStdString != "" {
		fmt.Printf("%s\n", outStdString)
	}
}

func DisableGoModReplaceDirectives(repoDir string, reposToReplace []string) error {
	goModPath := filepath.Join(repoDir, "go.mod")

	goModData, err := ioutil.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("could not read go.mod file at [%s]: %v", goModPath, err)
	}

	// compile expressions
	var replaceExpressions []*regexp.Regexp
	for _, repoToReplace := range reposToReplace {
		// look for shortcuts
		if newExpression, ok := knownRepos[repoToReplace]; ok {
			repoToReplace = newExpression
		}

		// test repo expression
		_, err := regexp.Compile(repoToReplace)
		if err != nil {
			return fmt.Errorf("could not compile regular expression from [%s], expressions must compile: %v", repoToReplace, err)
		}

		expressionStr := `(?m)^([ ]*replace[ ]+` + repoToReplace + `[ ]*)$`
		expression, err := regexp.Compile(expressionStr)
		if err != nil {
			return fmt.Errorf("could not format regular expressions from [%s]: target expression that failed [%s] error:%v", repoToReplace, expressionStr, err)
		}

		replaceExpressions = append(replaceExpressions, expression)
	}

	//replace in gomod
	for _, expression := range replaceExpressions {
		goModData = expression.ReplaceAll(goModData, []byte(`//$1`))
	}

	//write
	if err := ioutil.WriteFile(goModPath, goModData, os.ModePerm); err != nil {
		return fmt.Errorf("could not write go.mod file [%s]: %v", goModPath, err)
	}

	return nil
}

func EnableGoModReplaceDirectives(repoDir string, reposToReplace []string) error {
	goModPath := filepath.Join(repoDir, "go.mod")

	goModData, err := ioutil.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("could not read go.mod file at [%s]: %v", goModPath, err)
	}

	// compile expressions
	var replaceExpressions []*regexp.Regexp
	for _, repoToReplace := range reposToReplace {
		//look for shortcuts
		if newExpression, ok := knownRepos[repoToReplace]; ok {
			repoToReplace = newExpression
		}

		// test repo expression
		_, err := regexp.Compile(repoToReplace)
		if err != nil {
			return fmt.Errorf("could not compile regular expression from [%s], expressions must compile: %v", repoToReplace, err)
		}

		expressionStr := `(?m)^([ ]*)//([ ]*replace[ ]+` + repoToReplace + `[ ]*)$`
		expression, err := regexp.Compile(expressionStr)
		if err != nil {
			return fmt.Errorf("could not format regular expressions from [%s]: target expression that failed [%s] error:%v", repoToReplace, expressionStr, err)
		}

		replaceExpressions = append(replaceExpressions, expression)
	}

	//replace in gomod
	for _, expression := range replaceExpressions {
		goModData = expression.ReplaceAll(goModData, []byte(`$1$2`))
	}

	//write
	if err := ioutil.WriteFile(goModPath, goModData, os.ModePerm); err != nil {
		return fmt.Errorf("could not write go.mod file [%s]: %v", goModPath, err)
	}

	return nil
}

func HasGoModFile(repoDir string) bool {
	goModPath := filepath.Join(repoDir, "go.mod")

	info, err := os.Stat(goModPath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
