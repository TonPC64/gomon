package internal

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Builder interface
type Builder interface {
	Build() error
	Binary() string
	Errors() string
}

type builder struct {
	dir       string
	binary    string
	errors    string
	wd        string
	buildArgs []string
}

// NewBuilder creates new builder
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

	command := exec.Command(args[0], args[1:]...)

	output, err := command.CombinedOutput()
	if command.ProcessState.Success() {
		b.errors = ""
	} else {
		b.errors = string(output) + err.Error()
	}

	if len(b.errors) > 0 {
		return fmt.Errorf(b.errors)
	}

	return err
}
