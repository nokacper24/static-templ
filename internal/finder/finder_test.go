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
