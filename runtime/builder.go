package runtime

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Builder provides a binary builder
type Builder interface {
	// Build creates the temporal executable
	Build() error
	// Binary returns the reference to the runtime binary
	Binary() string
	// Errors returns any errors from the executable
	Errors() string
}

type builder struct {
	dir       string
	binary    string
	errors    string
	wd        string
	buildArgs []string
}

// New constructs a new Builder
func NewBuilder(dir string, bin string, wd string, buildArgs []string) Builder {
	if len(bin) == 0 {
		bin = "bin"
	}

	// does not work on Windows without the ".exe" extension
	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(bin, ".exe") { // check if it already has the .exe extension
			bin += ".exe"
		}
	}

	return &builder{dir: dir, binary: bin, wd: wd, buildArgs: buildArgs}
}

func (b *builder) Binary() string {
	return b.binary
}

func (b *builder) Errors() string {
	return b.errors
}

func (b *builder) Build() error {
	args := append([]string{"go", "build", "-o", filepath.Join(b.wd, b.binary)}, b.buildArgs...)

	var command *exec.Cmd
	command = exec.Command(args[0], args[1:]...)
	command.Dir = b.dir

	output, err := command.CombinedOutput()

	if command.ProcessState.Success() {
		b.errors = ""
	} else {
		b.errors = string(output)
	}

	if len(b.errors) > 0 {
		return fmt.Errorf(b.errors)
	}

	return err
}
