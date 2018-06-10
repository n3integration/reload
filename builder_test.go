package reload

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func Test_Builder_Build_Success(t *testing.T) {
	dir := filepath.Join("testdata", "build_success")
	bin := "build_success"
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not get working directory: %v", err)
	}

	builder := NewBuilder(dir, bin, wd, []string{})
	err = builder.Build()
	expect(t, err, nil)

	file, err := os.Open(filepath.Join(wd, bin))
	if err != nil {
		t.Fatalf("File has not been written: %v", err)
	}

	refute(t, file, nil)
}
