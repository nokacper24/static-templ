package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nokacper24/templ-static-generator/internal/finder"
	"github.com/nokacper24/templ-static-generator/internal/generator"
)

func main() {
	dirPath := "dist/"
	inputPath := "web/pages/"

	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Fatal("err creating dirs:", err)
	}

	modName, err := finder.FindModulePath()
	if err != nil {
		log.Fatal("Error finding module name:", err)
	}
	log.Println(modName)

	dirs, err := finder.FindDirsWithGoFiles(inputPath)
	if err != nil {
		log.Fatal("Error finding go dirs:", err)
	}
	log.Println(dirs)

	funcs, err := finder.FindFuncsToCall(dirs)
	if err != nil {
		log.Fatal("Error finding funcs:", err)
	}

	importsMap := map[string]bool{}
	for _, f := range funcs {
		// log.Println(f.FuncSign())
		// log.Println(f.HtmlFileName())
		importsMap[f.Location("")] = true
		fmt.Println(f.Location(""))
		fmt.Println(f.ToGenerate(inputPath, dirPath))
	}
	var importsSlice []string
	for imp := range importsMap {
		importsSlice = append(importsSlice, fmt.Sprintf("%s/%s", modName, imp))
	}

	err = generator.Generate(importsSlice, funcs)
	if err != nil {
		log.Fatal("err generating script", err)
	}

	// myfunc := funcs[0]
	// var funcName any = myfunc.FuncSign()
	// if v, ok := funcName.(templ.Component); ok {

	// 	fmt.Sprintf()
	// 	f, err := os.Create()
	// 	if err != nil {
	// 		log.Fatal("Error creating file:", err)
	// 	}
	// 	v.Render(context.Background(), f)
	// }

	// var funcName2 any = pages.Index()

	// ctx := context.Background()

	// p := "dist/test/file.html"

	// if err := os.MkdirAll(filepath.Dir(p), os.ModePerm); err != nil {
	// 	log.Fatal("error creating dirs:", err)
	// }

	// f, err := os.Create(p)
	// if err != nil {
	// 	log.Fatal("error creating file:", err)
	// }

	// pages.Index().Render(ctx,f)

	// f.Close()

}
