package analyzer

import (
	"go/ast"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "goimportalias",
	Doc:      "Checks all import aliases are consistent",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{ // filter needed nodes: visit only them
		(*ast.ImportSpec)(nil),
	}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		importStmt := node.(*ast.ImportSpec)

		if importStmt.Name == nil {
			return
		}
		alias := importStmt.Name.Name
		if alias == "" {
			return
		}
		aliasSlice := strings.Split(alias, "_")
		path := strings.ReplaceAll(importStmt.Path.Value, "\"", "")
		// replace all separators with `/` for normalization
		path = strings.ReplaceAll(path, "_", "/")
		path = strings.ReplaceAll(path, ".", "/")
		path = strings.ReplaceAll(path, "-", "")
		// omit the domain name in path
		pathSlice := strings.Split(path, "/")[1:]
		packageName := pathSlice[len(pathSlice)-1]

		if !checkVersion(aliasSlice[len(aliasSlice)-1], packageName) {
			pass.Reportf(node.Pos(), "version not specified in alias. path: %s alias: %s version %s", path, alias, packageName)
			return
		}
		if ok, lintErrMsg := checkAliasName(aliasSlice, pathSlice, pass); !ok {
			pass.Reportf(node.Pos(), lintErrMsg+" path: %s alias: %s", path, alias)
			return
		}
	})
	return nil, nil
}

// checkVersion checks that if package name starts with `v` it's included in alias name
func checkVersion(aliasLastWord string, packageName string) bool {
	if hasVPrefix := strings.HasPrefix(packageName, "v"); !hasVPrefix {
		return true
	}
	return aliasLastWord == packageName

}

// checkAliasName check consistency in alias name
func checkAliasName(aliasSlice []string, pathSlice []string, pass *analysis.Pass) (bool, string) {
	lastUsedWordIndex := -1
	for _, name := range aliasSlice {
		// we don't check version rule here
		if strings.HasPrefix(name, "v") || name == "" {
			continue
		}
		usedWordIndex := searchString(pathSlice, name)

		if usedWordIndex == len(pathSlice) {
			return false, "used words in alias most be present in path"
		}

		if usedWordIndex <= lastUsedWordIndex {
			return false, "order of words in alias should match words in path"
		}

		lastUsedWordIndex = usedWordIndex
	}

	if lastUsedWordIndex == -1 {
		return false, "at least one word from path must be present in alias"
	}

	return true, ""
}

func searchString(slice []string, word string) int {
	for pos, value := range slice {
		r, _ := regexp.Compile(word + "(s)?")
		if r.MatchString(value) {
			return pos
		}
	}
	return len(slice)
}
