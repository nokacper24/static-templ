package finder

import "testing"

func TestToGroupedFiles(t *testing.T) {
	filePaths := filePaths{
		"one/file.txt",
		"two/file.go",
		"theree/file_templ.go",
	}

	grouped := filePaths.ToGroupedFiles()

	if len(grouped.GoFiles) != 1 {
		t.Fatal("expected 1 go file")
	}
	if len(grouped.TemplGoFiles) != 1 {
		t.Fatal("expected 1 templ file")
	}
	if len(grouped.OtherFiles) != 1 {
		t.Fatal("expected 1 other file")
	}

	if grouped.GoFiles[0] != filePaths[1] {
		t.Fatal("expected go file to be...")
	}

	if grouped.TemplGoFiles[0] != filePaths[2] {
		t.Fatal("expected templ file to be...")
	}

	if grouped.OtherFiles[0] != filePaths[0] {
		t.Fatal("expected other file to be...")
	}
}

func TestRemoveTrailingSlash(t *testing.T) {
	str1 := "hello"
	strAfter1 := RemoveTrailingSlash(str1)
	if str1 != strAfter1 {
		t.Fatal("expected str1 to be equal to strAfter1")
	}

	str2 := ""
	strAfter2 := RemoveTrailingSlash(str2)
	if str2 != strAfter2 {
		t.Fatal("expected str2 to be equal to strAfter2")
	}

	str3 := "hello/"
	strAfter3 := RemoveTrailingSlash(str3)
	if strAfter3 != "hello" {
		t.Fatal("expected strAfter3 to be equal to hello")
	}

	str4 := "hello//"
	strAfter4 := RemoveTrailingSlash(str4)
	if strAfter4 != "hello/" {
		t.Fatal("expected strAfter4 to be equal to hello/")
	}

	str5 := "hel/lo/"
	strAfter5 := RemoveTrailingSlash(str5)
	if strAfter5 != "hel/lo" {
		t.Fatal("expected strAfter5 to be equal to hel/lo")
	}
}
