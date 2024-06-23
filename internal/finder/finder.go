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

type FunctionToCall struct {
	FileName     string // preserve the templ filename for the generated html one
	PackageName  string
	FunctionName string
	FilePath     string // used to determine import needed
}

// Finds path of the Go module the program is executed in.
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

// Finds paths to all files in the given directory and all its subdirecotries.
func findAllFiles(root string) (filePaths, error) {
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

// Goes throught the file paths provided and finds all exported fucntrions that take 0 parameters.
// Files provided must be valid Go source files.
func FindFunctionsInFiles(files filePaths) ([]FunctionToCall, error) {
	var funcs []FunctionToCall

	fileSt := token.NewFileSet()
	for _, path := range files {
		astFile, err := parser.ParseFile(fileSt, path, nil, parser.AllErrors)
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
						FunctionToCall{
							getFileNameWithoutExt(path),
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

// Returns path to the directory in which the function can be found.
func (f *FunctionToCall) DirPath() string {
	return filepath.Dir(f.FilePath)
}

// Returns the filename without path and extension.
func getFileNameWithoutExt(path string) string {
	fileName := filepath.Base(path) // Get the base filename with extension
	return fileName[:len(fileName)-9]
}

// Returns a string to be used as the name for HTML file generated from this component.
//
// Based on the original file name, with "-" replaced by "_", lowercased and .html added at the end.
func (f *FunctionToCall) HtmlFileName() string {
	noUnderscore := strings.ReplaceAll(f.FileName, "_", "-")
	lowered := strings.ToLower(noUnderscore)
	return fmt.Sprintf("%s.html", lowered)
}

type groupedFiles struct {
	TemplGoFiles filePaths // "_templ.go" files
	TemplFiles   filePaths // ".templ" files
	GoFiles      filePaths // ".go" files, excluding "_templ.go"
	OtherFiles   filePaths // other files
}

// Groups files provided into TemplGoFiles ("_templ.go"), TemplFiles (".templ"), GoFiles (other ".go" files) and OtherFiles.
func (f *filePaths) toGroupedFiles() *groupedFiles {
	var gf groupedFiles
	for _, fp := range *f {
		if fp[len(fp)-9:] == "_templ.go" {
			gf.TemplGoFiles = append(gf.TemplGoFiles, fp)
		} else if filepath.Ext(fp) == ".go" {
			gf.GoFiles = append(gf.GoFiles, fp)
		} else if filepath.Ext(fp) == ".templ" {
			gf.TemplFiles = append(gf.TemplFiles, fp)
		} else {
			gf.OtherFiles = append(gf.OtherFiles, fp)
		}
	}
	return &gf
}

// Finds paths to all files in the given directory and all its subdirecotries.
//
// Groups the files into groupedFiles type, includes TemplGoFiles ("_templ.go"), TemplFiles (".templ"), GoFiles (other ".go" files) and OtherFiles.
func FindFilesInDir(root string) (*groupedFiles, error) {
	allFiles, err := findAllFiles(root)
	if err != nil {
		return nil, err
	}
	return allFiles.toGroupedFiles(), nil
}

// Determines all imports needed to call provided functions.
func FindImports(funcs []FunctionToCall, modulePath string) []string {
	importsMap := map[string]bool{}
	for _, f := range funcs {
		importsMap[filepath.Dir(f.FilePath)] = true
	}
	var importsSlice []string
	for imp := range importsMap {
		importsSlice = append(importsSlice, fmt.Sprintf("%s/%s", modulePath, imp))
	}
	return importsSlice
}
