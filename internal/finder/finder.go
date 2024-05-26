package finder

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

func FindModulePath() (string, error) {
	modFile, err := os.ReadFile("go.mod")
	if err != nil {
		return "", err
	}

	mf, err := modfile.Parse("go.mod", modFile, nil)
	if err != nil {
		return "", err
	}

	return mf.Module.Mod.Path, nil
}

func FindDirsWithGoFiles(root string) ([]string, error) {
	goFilePaths := map[string]bool{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".go" {
			dirPath := filepath.Dir(path)
			goFilePaths[dirPath] = true
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var dirs []string
	for k := range goFilePaths {
		dirs = append(dirs, k)
	}
	return dirs, nil
}

func FindFuncsToCall(dirs []string) ([]FileToGenerate, error) {
	var funcs []FileToGenerate

	for _, path := range dirs {
		packages, err := parser.ParseDir(token.NewFileSet(), path, nil, parser.AllErrors)
		if err != nil {
			return nil, err
		}
		for packageName, v := range packages {
			for _, file := range v.Files {
				for _, d := range file.Decls {
					if fn, isFn := d.(*ast.FuncDecl); isFn {

						isExported := ast.IsExported(fn.Name.Name)
						hasParams := len(fn.Type.Params.List) != 0
						if isExported && !hasParams {
							functionName := fn.Name
							funcs = append(
								funcs,
								FileToGenerate{
									packageName,
									functionName.Name,
									path,
								})
						}
					}
				}
			}
		}

	}
	return funcs, nil
}

type FileToGenerate struct {
	PackageName  string
	FunctionName string
	path         string
}

func (f *FileToGenerate) FuncSign() string {
	return fmt.Sprintf("%s.%s()", f.PackageName, f.FunctionName)
}

func (f *FileToGenerate) HtmlFileName() string {
	noUnderscore := strings.ReplaceAll(f.FunctionName, "_", "-")
	lowered := strings.ToLower(noUnderscore)
	return fmt.Sprintf("%s.html", lowered)
}

// file location skipping the root
func (f *FileToGenerate) Location(root string) string {
	return f.path[len(root):]
}

func (f *FileToGenerate) ToGenerate(root string, prefix string) string {
	noRoot := f.path[len(root):]
	return fmt.Sprintf("%s%s/%s", prefix, noRoot, f.HtmlFileName())
}
