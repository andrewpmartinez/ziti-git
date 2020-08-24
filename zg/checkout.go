package zg

import (
	"fmt"
	"golang.org/x/mod/modfile"
	"io/ioutil"
	"path/filepath"
)

type VersionInfo struct {
}

func GetGoModInfo(repoDir string) (*modfile.File, error) {
	goModPath := filepath.Join(repoDir, "go.mod")

	goModData, err := ioutil.ReadFile(goModPath)
	if err != nil {
		return nil, fmt.Errorf("could not read go.mod file at [%s]: %v", goModPath, err)
	}

	return modfile.Parse(goModPath, goModData, nil)
}
