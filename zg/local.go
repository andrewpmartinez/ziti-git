package zg

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

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

func DisableGoModReplaceDirectives(repoDir string, reposToReplace []string) error {
	goModPath := filepath.Join(repoDir, "go.mod")

	goModData, err := ioutil.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("could not read go.mod file at [%s]: %v", goModPath, err)
	}

	// compile expressions
	var replaceExpressions []*regexp.Regexp
	for _, repoToReplace := range reposToReplace {
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
