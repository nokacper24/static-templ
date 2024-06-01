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

type filePaths []string

func FindAllFiles(root string) (filePaths, error) {
	var paths filePaths
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			paths = append(paths, path)

		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return paths, nil
}

func FindFuncsToCall(files filePaths) ([]FileToGenerate, error) {
	var funcs []FileToGenerate

	for _, path := range files {
		astFile, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.AllErrors)
		if err != nil {
			return nil, err
		}
		packageName := astFile.Name.Name
		for _, d := range astFile.Decls {
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
	return funcs, nil
}

type FileToGenerate struct {
	PackageName  string
	FunctionName string
	filePath     string
}

func (f *FileToGenerate) DirPath() string {
	return filepath.Dir(f.filePath)
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
	return f.filePath[len(root):]
}

func (f *FileToGenerate) ToGenerate(root string, prefix string) string {
	return fmt.Sprint(
		strings.Replace(f.DirPath(), root, prefix, 1),
		"/",
		f.HtmlFileName())
}

type groupedFiles struct {
	TemplGoFiles filePaths
	GoFiles      filePaths
	OtherFiles   filePaths
}

func (f *filePaths) ToGroupedFiles() *groupedFiles {
	var gf groupedFiles

	for _, fp := range *f {
		if fp[len(fp)-9:] == "_templ.go" {
			gf.TemplGoFiles = append(gf.TemplGoFiles, fp)
		} else if filepath.Ext(fp) == ".go" {
			gf.GoFiles = append(gf.GoFiles, fp)
		} else {
			gf.OtherFiles = append(gf.OtherFiles, fp)
		}
	}

	return &gf
}
